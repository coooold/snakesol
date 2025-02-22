// Package network 处理WebSocket连接和消息通信
package network

import (
	"encoding/json"
	"log"
	"net/http"

	"snakesol/internal/game"

	"github.com/gorilla/websocket"
)

// WSServer 处理WebSocket连接的服务器
type WSServer struct {
	game     *game.GameState
	upgrader websocket.Upgrader
}

// Message WebSocket消息结构
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewWSServer 创建一个新的WebSocket服务器
func NewWSServer(game *game.GameState) *WSServer {
	return &WSServer{
		game: game,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// HandleConnection 处理新的WebSocket连接
func (s *WSServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级WebSocket连接失败:", err)
		return
	}

	// 创建新玩家的蛇
	snake := game.CreateSnake(false)
	wsConn := &WSConnection{conn: conn}
	snake.Conn = wsConn

	// 将蛇添加到游戏状态
	s.game.AddSnake(snake)

	// 处理玩家输入
	go s.handlePlayerInput(snake)
}

// handlePlayerInput 处理玩家的输入消息
func (s *WSServer) handlePlayerInput(snake *game.Snake) {
	wsConn := snake.Conn.(*WSConnection)

	for {
		var msg Message
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			// 从游戏状态中移除蛇，并将其转换为苹果
			s.game.RemoveSnake(snake.ID)
			// 关闭WebSocket连接
			wsConn.Close()
			return
		}

		if msg.Type == "direction" {
			var dir game.Direction
			if err := json.Unmarshal(msg.Payload, &dir); err == nil {
				s.game.UpdateSnakeDirection(snake.ID, dir)
			}
		}
	}
}

// WSConnection 包装websocket.Conn以实现game.Connection接口
type WSConnection struct {
	conn *websocket.Conn
}

// WriteJSON 实现game.Connection接口
func (c *WSConnection) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

// ReadJSON 实现game.Connection接口
func (c *WSConnection) ReadJSON(v interface{}) error {
	return c.conn.ReadJSON(v)
}

// Close 实现game.Connection接口
func (c *WSConnection) Close() error {
	return c.conn.Close()
}
