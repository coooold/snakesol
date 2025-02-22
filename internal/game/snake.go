package game

import (
	"math/rand"
	"time"
)

// Snake 表示一条蛇
type Snake struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Color       string          `json:"color"`
	IsAI        bool            `json:"isAI"`
	X           int             `json:"x"`
	Y           int             `json:"y"`
	Direction   Direction       `json:"direction"`
	Body        []Position      `json:"body"`
	Dead        bool            `json:"dead"`
	Conn        Connection      `json:"-"`
	Personality PersonalityType `json:"personality"`
}

// CreateSnake 创建一条新蛇
func CreateSnake(isAI bool) *Snake {
	config := DefaultConfig()
	x := rand.Intn(config.Cols)
	y := rand.Intn(config.Rows)
	dirs := []Direction{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	dir := dirs[rand.Intn(len(dirs))]

	personality := RandomPersonality()

	snake := &Snake{
		ID:          generateID(),
		Name:        getRandName(),
		Color:       PersonalityColor[personality], // 使用性格对应的颜色
		IsAI:        isAI,
		X:           x,
		Y:           y,
		Direction:   dir,
		Body:        make([]Position, 0),
		Personality: personality,
	}

	snakeLen := config.InitialSnakeLength
	if !isAI {
		snakeLen = 20
	}

	// 初始化蛇身
	for i := 0; i < snakeLen; i++ {
		snake.Body = append(snake.Body, Position{
			X: (x - dir.X*i + config.Cols) % config.Cols,
			Y: (y - dir.Y*i + config.Rows) % config.Rows,
		})
	}

	return snake
}

// generateID 生成唯一的蛇ID
func generateID() string {
	return time.Now().Format("20060102150405") +
		string(rune(rand.Intn(26)+'A'))
}
