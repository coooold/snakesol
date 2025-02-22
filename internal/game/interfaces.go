// Package game 定义游戏领域的核心接口
package game

// Connection 定义了与客户端通信的接口
type Connection interface {
    WriteJSON(v interface{}) error
    ReadJSON(v interface{}) error
    Close() error
}

// EventEmitter 定义了游戏事件发射器的接口
type EventEmitter interface {
    Emit(eventType string, data interface{})
    On(eventType string, handler func(data interface{}))
    Off(eventType string, handler func(data interface{}))
}