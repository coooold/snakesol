package game

import (
	"math/rand"
	"sync"
	"time"
)

// GameState 实现了游戏状态管理
type GameState struct {
	snakes map[string]*Snake
	apples []AppleInfo
	mu     sync.Mutex
	config *GameConfig
}

type AppleInfo struct {
	Position  Position
	CreatedAt time.Time
}

// NewGameState 创建一个新的游戏状态
func NewGameState() *GameState {
	gs := &GameState{
		snakes: make(map[string]*Snake),
		apples: make([]AppleInfo, 0),
		config: DefaultConfig(),
	}

	// 初始化时添加AI蛇
	for i := 0; i < gs.config.InitialAICount; i++ {
		// 创建新的AI蛇
		snake := CreateSnake(true)

		// 检查生成位置是否与其他蛇重叠
		isValidPosition := true
		for _, existingSnake := range gs.snakes {
			if !existingSnake.Dead {
				// 检查头部位置
				if snake.X == existingSnake.X && snake.Y == existingSnake.Y {
					isValidPosition = false
					break
				}
				// 检查身体位置
				for _, segment := range existingSnake.Body {
					if snake.X == segment.X && snake.Y == segment.Y {
						isValidPosition = false
						break
					}
				}
			}
		}

		// 如果位置有效，则添加到游戏中
		if isValidPosition {
			gs.snakes[snake.ID] = snake
		} else {
			// 如果位置无效，重试这一次
			i--
		}
	}

	go gs.spawnAISnakes()
	go gs.spawnApples() // 启动苹果生成器
	return gs
}

// AddSnake 添加一条蛇到游戏中
func (gs *GameState) AddSnake(snake *Snake) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.snakes[snake.ID] = snake
	gs.broadcastState()
}

// RemoveSnake 从游戏中移除一条蛇
func (gs *GameState) RemoveSnake(id string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	// 如果蛇还存在，则将其转换为苹果
	if snake, ok := gs.snakes[id]; ok {
		gs.snakeToApples(snake)
	}
	gs.broadcastState()
}

// UpdateSnakeDirection 更新蛇的移动方向
func (gs *GameState) UpdateSnakeDirection(id string, dir Direction) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if snake, ok := gs.snakes[id]; ok {
		snake.Direction = dir
	}
}

// UpdateGame 更新游戏状态
func (gs *GameState) UpdateGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// 检查并移除超过20秒的苹果
	now := time.Now()
	var validApples []AppleInfo
	for _, apple := range gs.apples {
		if now.Sub(apple.CreatedAt).Seconds() < float64(gs.config.AppleLifetime) {
			validApples = append(validApples, apple)
		}
	}
	if len(validApples) != len(gs.apples) {
		gs.apples = validApples
		gs.broadcastState()
	}

	// 更新AI蛇的方向
	for _, snake := range gs.snakes {
		if snake.IsAI && !snake.Dead {
			ai := NewAIController(snake, gs.config)
			snake.Direction = ai.DecideNextMove(gs)
		}
	}

	// 更新所有蛇的位置
	for _, snake := range gs.snakes {
		if snake.Dead {
			continue
		}

		// 移动蛇身
		snake.Body = append([]Position{{X: snake.X, Y: snake.Y}}, snake.Body...)
		snake.X = (snake.X + snake.Direction.X + gs.config.Cols) % gs.config.Cols
		snake.Y = (snake.Y + snake.Direction.Y + gs.config.Rows) % gs.config.Rows

		// 检查是否吃到苹果
		for i, apple := range gs.apples {
			if apple.Position.X == snake.X && apple.Position.Y == snake.Y {
				gs.apples = append(gs.apples[:i], gs.apples[i+1:]...)
				goto skipTail
			}
		}
		if len(snake.Body) > 0 {
			snake.Body = snake.Body[:len(snake.Body)-1]
		}
	skipTail:

		// 检查自身碰撞
		for _, segment := range snake.Body {
			if segment.X == snake.X && segment.Y == snake.Y {
				// 处理蛇的死亡，转换为苹果
				gs.snakeToApples(snake)
				goto nextSnake
			}
		}

		// 检查与其他蛇的碰撞
		for _, other := range gs.snakes {
			if other == snake || other.Dead {
				continue
			}
			for _, segment := range other.Body {
				if segment.X == snake.X && segment.Y == snake.Y {
					// 处理蛇的死亡，转换为苹果
					gs.snakeToApples(snake)
					goto nextSnake
				}
			}
		}
	nextSnake:
	}

	gs.broadcastState()
}

// broadcastState 向所有玩家广播游戏状态
func (gs *GameState) broadcastState() {
	state := struct {
		Snakes      map[string]*Snake `json:"snakes"`
		Apples      []Position        `json:"apples"`
		Config      *GameConfig       `json:"config"`
		DeadSnakeID string            `json:"deadSnakeId,omitempty"`
	}{
		Snakes: gs.snakes,
		Apples: func() []Position {
			positions := make([]Position, len(gs.apples))
			for i, apple := range gs.apples {
				positions[i] = apple.Position
			}
			return positions
		}(),
		Config: gs.config,
	}

	for _, snake := range gs.snakes {
		if !snake.IsAI && snake.Conn != nil {
			snake.Conn.WriteJSON(state)
		}
	}
}

// snakeToApples 将死亡的蛇转换为苹果
func (gs *GameState) snakeToApples(snake *Snake) {
	// 设置蛇的死亡状态
	snake.Dead = true

	// 广播死亡事件
	deathEvent := struct {
		Snakes      map[string]*Snake `json:"snakes"`
		Apples      []Position        `json:"apples"`
		Config      *GameConfig       `json:"config"`
		DeadSnakeID string            `json:"deadSnakeId"`
	}{
		Snakes: gs.snakes,
		Apples: func() []Position {
			positions := make([]Position, len(gs.apples))
			for i, apple := range gs.apples {
				positions[i] = apple.Position
			}
			return positions
		}(),
		Config:      gs.config,
		DeadSnakeID: snake.ID,
	}

	// 向所有玩家广播死亡事件
	for _, s := range gs.snakes {
		if !s.IsAI && s.Conn != nil {
			s.Conn.WriteJSON(deathEvent)
		}
	}

	// 在蛇身体的每个位置生成苹果
	for _, segment := range snake.Body {
		gs.apples = append(gs.apples, AppleInfo{
			Position:  segment,
			CreatedAt: time.Now(),
		})
	}
	// 从游戏中移除死亡的蛇
	delete(gs.snakes, snake.ID)
}

// spawnAISnakes 定期生成AI控制的蛇
func (gs *GameState) spawnAISnakes() {
	ticker := time.NewTicker(time.Duration(gs.config.AISpawnInterval) * time.Second) // 固定10秒生成一个AI
	for range ticker.C {
		gs.mu.Lock()
		// 检查当前AI数量
		aiCount := 0
		for _, s := range gs.snakes {
			if s.IsAI && !s.Dead {
				aiCount++
			}
		}

		// 限制场景中的AI数量
		if aiCount < gs.config.MaxAICount {
			// 创建新的AI蛇
			snake := CreateSnake(true)

			// 检查生成位置是否与其他蛇重叠
			isValidPosition := true
			for _, existingSnake := range gs.snakes {
				if !existingSnake.Dead {
					// 检查头部位置
					if snake.X == existingSnake.X && snake.Y == existingSnake.Y {
						isValidPosition = false
						break
					}
					// 检查身体位置
					for _, segment := range existingSnake.Body {
						if snake.X == segment.X && snake.Y == segment.Y {
							isValidPosition = false
							break
						}
					}
				}
			}

			// 如果位置有效，则添加到游戏中
			if isValidPosition {
				gs.snakes[snake.ID] = snake
				gs.broadcastState()
			}
		}
		gs.mu.Unlock()
	}
}

// spawnApples 定期在随机位置生成苹果
func (gs *GameState) spawnApples() {
	ticker := time.NewTicker(time.Duration(gs.config.AppleSpawnInterval) * time.Second)
	for range ticker.C {
		gs.mu.Lock()
		// 生成随机位置
		x := rand.Intn(gs.config.Cols)
		y := rand.Intn(gs.config.Rows)

		// 检查该位置是否已经有苹果或蛇
		isValidPosition := true
		// 检查是否与现有苹果重叠
		for _, apple := range gs.apples {
			if apple.Position.X == x && apple.Position.Y == y {
				isValidPosition = false
				break
			}
		}
		// 检查是否与蛇重叠
		if isValidPosition {
			for _, snake := range gs.snakes {
				if !snake.Dead {
					if snake.X == x && snake.Y == y {
						isValidPosition = false
						break
					}
					for _, segment := range snake.Body {
						if segment.X == x && segment.Y == y {
							isValidPosition = false
							break
						}
					}
				}
			}
		}

		// 如果位置有效，添加苹果
		if isValidPosition {
			gs.apples = append(gs.apples, AppleInfo{
				Position:  Position{X: x, Y: y},
				CreatedAt: time.Now(),
			})
			gs.broadcastState()
		}
		gs.mu.Unlock()
	}
}
