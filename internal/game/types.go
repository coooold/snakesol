package game

// Direction 表示移动方向
type Direction struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Position 表示坐标位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GameConfig 游戏配置
type GameConfig struct {
	Cols               int `json:"cols"`
	Rows               int `json:"rows"`
	InitialSnakeLength int
	AISpawnInterval    int
	UpdateInterval     int
	InitialAICount     int
	MaxAICount         int
	AppleSpawnInterval int
	AppleLifetime      int
}
