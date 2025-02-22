package game

import (
	"fmt"
	"math"
)

// AIController 处理AI蛇的决策逻辑
type AIController struct {
	snake  *Snake
	config *GameConfig
}

// NewAIController 创建新的AI控制器
func NewAIController(snake *Snake, config *GameConfig) *AIController {
	return &AIController{
		snake:  snake,
		config: config,
	}
}

// DecideNextMove 决定AI蛇的下一步移动方向
func (ai *AIController) DecideNextMove(gameState *GameState) Direction {
	// 获取当前可用的移动方向
	availableDirections := ai.getAvailableDirections(gameState)
	if len(availableDirections) == 0 {
		return ai.snake.Direction // 如果没有安全的方向，保持当前方向
	}

	// 评估每个方向的得分
	bestDirection := ai.snake.Direction
	bestScore := float64(-1000000)

	for _, dir := range availableDirections {
		score := ai.evaluateDirection(dir, gameState)
		if score > bestScore {
			bestScore = score
			bestDirection = dir
		}
	}

	return bestDirection
}

// getAvailableDirections 获取安全的移动方向
func (ai *AIController) getAvailableDirections(gameState *GameState) []Direction {
	directions := []Direction{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	safeDirections := make([]Direction, 0)

	// 获取当前方向的反方向
	oppositeDir := Direction{-ai.snake.Direction.X, -ai.snake.Direction.Y}

	for _, dir := range directions {
		// 检查是否是反方向
		if dir == oppositeDir {
			continue
		}

		nextX := (ai.snake.X + dir.X + ai.config.Cols) % ai.config.Cols
		nextY := (ai.snake.Y + dir.Y + ai.config.Rows) % ai.config.Rows

		// 检查是否会撞到自己
		safe := true
		for _, segment := range ai.snake.Body {
			if nextX == segment.X && nextY == segment.Y {
				safe = false
				break
			}
		}
		if !safe {
			continue
		}

		// 检查是否会撞到其他蛇
		for _, otherSnake := range gameState.snakes {
			if otherSnake.Dead || otherSnake == ai.snake {
				continue
			}
			// 检查头部碰撞
			if nextX == otherSnake.X && nextY == otherSnake.Y {
				safe = false
				break
			}
			// 检查身体碰撞
			for _, segment := range otherSnake.Body {
				if nextX == segment.X && nextY == segment.Y {
					safe = false
					break
				}
			}
		}

		if safe {
			safeDirections = append(safeDirections, dir)
		}
	}

	return safeDirections
}

// evaluateDirection 评估某个方向的得分
func (ai *AIController) evaluateDirection(dir Direction, gameState *GameState) float64 {
	nextX := (ai.snake.X + dir.X + ai.config.Cols) % ai.config.Cols
	nextY := (ai.snake.Y + dir.Y + ai.config.Rows) % ai.config.Rows

	// 获取性格权重并根据局势动态调整
	baseWeights := GetPersonalityWeights(ai.snake.Personality)
	weights := ai.adjustWeightsByGameState(baseWeights, gameState)

	// 获取视野范围内的信息
	viewInfo := ai.snake.GetViewInfo(gameState)

	// 基础得分
	score := 0.0

	// 1. 食物评分 - 根据视野内的食物评估
	minFoodDist := float64(ViewWidth + ViewHeight)
	for _, food := range viewInfo.Food {
		dist := ai.distance(nextX, nextY, food.X, food.Y)
		if dist < minFoodDist {
			minFoodDist = dist
		}
	}
	score += (float64(ViewWidth+ViewHeight) - minFoodDist) * 10 * weights.FoodWeight

	// 2. 空间评分 - 根据视野内的空间评估
	spaceScore := ai.evaluateSpace(nextX, nextY, gameState)
	score += spaceScore * 5 * weights.SpaceWeight

	// 3. 生存评分 - 根据视野内的威胁评估
	survivalScore := 0.0
	for _, otherSnake := range viewInfo.Snakes {
		// 与其他蛇的距离
		dist := ai.distance(nextX, nextY, otherSnake.X, otherSnake.Y)
		if dist < 2 { // 如果太近，根据性格降低得分
			survivalScore -= (2 - dist) * 100
		}
	}
	score += survivalScore * weights.SurvivalWeight

	// 4. 攻击评分 - 根据视野内的攻击机会评估
	attackScore := 0.0
	if len(ai.snake.Body) > 10 {
		for _, otherSnake := range viewInfo.Snakes {
			if len(otherSnake.Body) < len(ai.snake.Body) {
				dist := ai.distance(nextX, nextY, otherSnake.X, otherSnake.Y)
				if dist < 5 { // 根据性格评估攻击价值
					attackScore += (5 - dist) * 20
				}
			}
		}
	}
	score += attackScore * weights.AttackWeight

	// 5. 机动性评分 - 特别适用于游击型
	if weights.MobilityWeight > 0 {
		mobilityScore := float64(len(ai.getAvailableDirections(gameState)))
		score += mobilityScore * 10 * weights.MobilityWeight
	}

	// 6. 陷阱评分 - 特别适用于陷阱型
	if weights.TrapWeight > 0 {
		trapScore := ai.evaluateTrapPotential(nextX, nextY, gameState)
		score += trapScore * weights.TrapWeight
	}

	// 7. 协作评分 - 特别适用于协作型
	if weights.CooperateWeight > 0 {
		cooperateScore := ai.evaluateCooperation(nextX, nextY, gameState)
		score += cooperateScore * weights.CooperateWeight
	}

	return score
}

// evaluateSpace 评估某个位置的可用空间
func (ai *AIController) evaluateSpace(x, y int, gameState *GameState) float64 {
	visited := make(map[string]bool)
	space := ai.floodFill(x, y, gameState, visited)
	return float64(space)
}

// floodFill 使用泛洪算法计算可用空间
func (ai *AIController) floodFill(x, y int, gameState *GameState, visited map[string]bool) int {
	key := fmt.Sprintf("%d,%d", x, y)
	if visited[key] {
		return 0
	}

	// 检查是否是障碍物
	for _, snake := range gameState.snakes {
		if snake.Dead {
			continue
		}
		if x == snake.X && y == snake.Y {
			return 0
		}
		for _, segment := range snake.Body {
			if x == segment.X && y == segment.Y {
				return 0
			}
		}
	}

	visited[key] = true
	space := 1

	// 限制搜索深度以提高性能
	if len(visited) > 50 {
		return space
	}

	// 递归搜索四个方向
	directions := []Direction{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for _, dir := range directions {
		nextX := (x + dir.X + ai.config.Cols) % ai.config.Cols
		nextY := (y + dir.Y + ai.config.Rows) % ai.config.Rows
		space += ai.floodFill(nextX, nextY, gameState, visited)
	}

	return space
}

// evaluateTrapPotential 评估某个位置设置陷阱的潜力
func (ai *AIController) evaluateTrapPotential(x, y int, gameState *GameState) float64 {
    // 基础陷阱得分
    trapScore := 0.0

    // 检查周围的通道数量
    passages := 0
    directions := []Direction{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
    for _, dir := range directions {
        nextX := (x + dir.X + ai.config.Cols) % ai.config.Cols
        nextY := (y + dir.Y + ai.config.Rows) % ai.config.Rows
        
        // 检查该方向是否有障碍物
        hasObstacle := false
        for _, snake := range gameState.snakes {
            if snake.Dead {
                continue
            }
            if nextX == snake.X && nextY == snake.Y {
                hasObstacle = true
                break
            }
            for _, segment := range snake.Body {
                if nextX == segment.X && nextY == segment.Y {
                    hasObstacle = true
                    break
                }
            }
        }
        
        if !hasObstacle {
            passages++
        }
    }

    // 理想的陷阱位置应该只有1-2个通道
    if passages >= 1 && passages <= 2 {
        trapScore += float64(50 * (3 - passages))

        // 检查附近是否有其他蛇
        for _, otherSnake := range gameState.snakes {
            if otherSnake.Dead || otherSnake == ai.snake {
                continue
            }
            
            dist := ai.distance(x, y, otherSnake.X, otherSnake.Y)
            if dist < 5 {
                // 如果附近有较小的蛇，增加陷阱得分
                if len(otherSnake.Body) < len(ai.snake.Body) {
                    trapScore += (5 - dist) * 30
                } else {
                    // 如果附近有较大的蛇，降低陷阱得分
                    trapScore -= (5 - dist) * 20
                }
            }
        }

        // 确保有逃生路线
        escapeSpace := ai.evaluateSpace(x, y, gameState)
        if escapeSpace > 10 {
            trapScore += float64(escapeSpace) * 2
        } else {
            trapScore -= (10 - escapeSpace) * 10
        }
    }

    return math.Max(0, trapScore)
}

// distance 计算两点之间的曼哈顿距离
func (ai *AIController) distance(x1, y1, x2, y2 int) float64 {
	dx := math.Abs(float64(x1 - x2))
	dy := math.Abs(float64(y1 - y2))
	// 考虑游戏场地是环形的情况
	dx = math.Min(dx, float64(ai.config.Cols)-dx)
	dy = math.Min(dy, float64(ai.config.Rows)-dy)
	return dx + dy
}

// evaluateCooperation 评估协作行为的得分
func (ai *AIController) evaluateCooperation(x, y int, gameState *GameState) float64 {
	cooperateScore := 0.0

	// 获取视野范围内的信息
	viewInfo := ai.snake.GetViewInfo(gameState)

	// 寻找同为协作型的其他蛇
	cooperativeSnakes := make([]*Snake, 0)
	for _, otherSnake := range viewInfo.Snakes {
		if otherSnake.Personality == Cooperative {
			cooperativeSnakes = append(cooperativeSnakes, otherSnake)
		}
	}

	// 如果没有其他协作型蛇，返回0分
	if len(cooperativeSnakes) == 0 {
		return 0
	}

	// 寻找目标蛇（非协作型且体型较小的蛇）
	var targetSnake *Snake
	for _, snake := range viewInfo.Snakes {
		if snake.Personality != Cooperative && len(snake.Body) < len(ai.snake.Body) {
			if targetSnake == nil || len(snake.Body) < len(targetSnake.Body) {
				targetSnake = snake
			}
		}
	}

	// 如果找到目标蛇，实施围堵策略
	if targetSnake != nil {
		// 计算目标蛇到各个协作蛇的平均距离
		totalDist := 0.0
		for _, cooperativeSnake := range cooperativeSnakes {
			totalDist += ai.distance(cooperativeSnake.X, cooperativeSnake.Y, targetSnake.X, targetSnake.Y)
		}
		avgDist := totalDist / float64(len(cooperativeSnakes))

		// 计算当前位置到目标蛇的距离
		currentDist := ai.distance(x, y, targetSnake.X, targetSnake.Y)

		// 如果当前距离小于平均距离，增加得分
		if currentDist < avgDist {
			cooperateScore += (avgDist - currentDist) * 30
		}

		// 根据与其他协作蛇的相对位置调整得分
		for _, cooperativeSnake := range cooperativeSnakes {
			// 计算当前蛇和协作蛇与目标的夹角
			angle := ai.calculateAngle(x, y, cooperativeSnake.X, cooperativeSnake.Y, targetSnake.X, targetSnake.Y)
			
			// 如果夹角接近90度，说明形成了较好的包围态势
			if angle >= 45 && angle <= 135 {
				cooperateScore += 50
			}
		}
	}

	return cooperateScore
}

// calculateAngle 计算三点形成的夹角（度数）
func (ai *AIController) calculateAngle(x1, y1, x2, y2, x3, y3 int) float64 {
	// 将坐标转换为相对于x3,y3的向量
	vector1X := float64(x1 - x3)
	vector1Y := float64(y1 - y3)
	vector2X := float64(x2 - x3)
	vector2Y := float64(y2 - y3)

	// 计算向量的点积
	dotProduct := vector1X*vector2X + vector1Y*vector2Y

	// 计算向量的模
	magnitude1 := math.Sqrt(vector1X*vector1X + vector1Y*vector1Y)
	magnitude2 := math.Sqrt(vector2X*vector2X + vector2Y*vector2Y)

	// 计算夹角（弧度）
	angle := math.Acos(dotProduct / (magnitude1 * magnitude2))

	// 转换为度数
	return angle * 180 / math.Pi
}

// adjustWeightsByGameState 根据游戏局势动态调整权重
func (ai *AIController) adjustWeightsByGameState(baseWeights PersonalityWeights, gameState *GameState) PersonalityWeights {
	adjustedWeights := baseWeights

	// 获取当前蛇的状态
	snakeLength := len(ai.snake.Body)
	viewInfo := ai.snake.GetViewInfo(gameState)

	// 1. 根据蛇的长度调整权重
	if snakeLength < 10 {
		// 当蛇较小时，提高生存权重，降低攻击权重
		adjustedWeights.SurvivalWeight *= 1.5
		adjustedWeights.AttackWeight *= 0.5
		adjustedWeights.FoodWeight *= 1.3
	} else if snakeLength > 30 {
		// 当蛇较大时，提高攻击权重
		adjustedWeights.AttackWeight *= 1.3
		adjustedWeights.SurvivalWeight *= 0.8
	}

	// 2. 根据周围环境调整权重
	nearbySnakes := 0
	largerSnakes := 0
	for _, snake := range viewInfo.Snakes {
		dist := ai.distance(ai.snake.X, ai.snake.Y, snake.X, snake.Y)
		if dist < 10 {
			nearbySnakes++
			if len(snake.Body) > snakeLength {
				largerSnakes++
			}
		}
	}

	// 3. 根据威胁程度调整权重
	if largerSnakes > 0 {
		// 周围有较大的蛇时，提高生存和空间权重
		adjustedWeights.SurvivalWeight *= 1.5
		adjustedWeights.SpaceWeight *= 1.3
		adjustedWeights.AttackWeight *= 0.6
	}

	// 4. 根据食物分布调整权重
	if len(viewInfo.Food) > 3 {
		// 周围食物较多时，适当提高觅食权重
		adjustedWeights.FoodWeight *= 1.2
	}

	// 5. 协作型特殊调整
	if ai.snake.Personality == Cooperative {
		cooperativeNearby := false
		for _, snake := range viewInfo.Snakes {
			if snake.Personality == Cooperative {
				cooperativeNearby = true
				break
			}
		}
		if cooperativeNearby {
			// 附近有其他协作型蛇时，强化协作权重
			adjustedWeights.CooperateWeight *= 1.5
		}
	}

	return adjustedWeights
}