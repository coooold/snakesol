package game

// GetViewInfo 获取AI蛇的视野范围内的信息
func (s *Snake) GetViewInfo(gameState *GameState) ViewInfo {
	// 创建视野信息结构，增加威胁等级评估
	view := ViewInfo{
		Center:    Position{X: s.X, Y: s.Y},
		Food:      make([]Position, 0),
		Snakes:    make([]*Snake, 0),
		Obstacles: make([]Position, 0),
	}

	// 计算视野范围
	minX := s.X - ViewWidth/2
	maxX := s.X + ViewWidth/2
	minY := s.Y - ViewHeight/2
	maxY := s.Y + ViewHeight/2

	// 获取视野范围内的食物
	for _, food := range gameState.apples {
		if isInView(food.Position.X, food.Position.Y, minX, maxX, minY, maxY) {
			view.Food = append(view.Food, food.Position)
		}
	}

	// 获取视野范围内的其他蛇
	for _, snake := range gameState.snakes {
		if snake.ID != s.ID && !snake.Dead {
			if isInView(snake.X, snake.Y, minX, maxX, minY, maxY) {
				view.Snakes = append(view.Snakes, snake)
			}
			// 将其他蛇的身体部分作为障碍物
			for _, pos := range snake.Body {
				if isInView(pos.X, pos.Y, minX, maxX, minY, maxY) {
					view.Obstacles = append(view.Obstacles, pos)
				}
			}
		}
	}

	return view
}

// isInView 判断一个位置是否在视野范围内
func isInView(x, y, minX, maxX, minY, maxY int) bool {
	// 处理地图边界循环的情况
	config := DefaultConfig()

	// 将坐标转换到合法范围内
	x = (x + config.Cols) % config.Cols
	y = (y + config.Rows) % config.Rows
	minX = (minX + config.Cols) % config.Cols
	maxX = (maxX + config.Cols) % config.Cols
	minY = (minY + config.Rows) % config.Rows
	maxY = (maxY + config.Rows) % config.Rows

	// 处理跨越边界的情况
	if minX > maxX {
		return x >= minX || x <= maxX
	}
	if minY > maxY {
		return y >= minY || y <= maxY
	}

	return x >= minX && x <= maxX && y >= minY && y <= maxY
}
