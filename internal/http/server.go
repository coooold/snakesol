// Package http 提供HTTP和WebSocket服务
package http

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"snakesol/internal/game"
	"snakesol/internal/network"
)

// Config HTTP服务器配置
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// Server HTTP服务器
type Server struct {
	config    *Config
	gameState *game.GameState
	wsServer  *network.WSServer
	staticFS  embed.FS
}

// NewServer 创建HTTP服务器
func NewServer(config *Config, gameState *game.GameState, staticFS embed.FS) *Server {
	return &Server{
		config:    config,
		gameState: gameState,
		wsServer:  network.NewWSServer(gameState),
		staticFS:  staticFS,
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	// 获取嵌入的静态文件系统
	subFS, err := fs.Sub(s.staticFS, "client")
	if err != nil {
		return err
	}

	// 设置静态文件服务
	fs := http.FileServer(http.FS(subFS))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

	// 设置WebSocket路由
	http.HandleFunc("/ws", s.wsServer.HandleConnection)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         s.config.Addr,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	log.Printf("游戏服务器启动在 %s", s.config.Addr)
	return server.ListenAndServe()
}
