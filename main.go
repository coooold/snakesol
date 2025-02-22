package main

import (
	"embed"
	"flag"
	"log"
	"math/rand"
	"time"

	"snakesol/internal/game"
	"snakesol/internal/http"
)

//go:generate go run github.com/markbates/pkger/cmd/pkger -o server

//go:embed client/*
var staticFiles embed.FS

func main() {
	// 解析命令行参数
	port := flag.String("port", "8080", "服务器监听端口")
	flag.Parse()

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 创建游戏状态
	gameState := game.NewGameState()

	// 启动游戏循环
	go func() {
		ticker := time.NewTicker(time.Duration(game.DefaultConfig().UpdateInterval) * time.Millisecond)
		for range ticker.C {
			gameState.UpdateGame()
		}
	}()

	// 创建HTTP服务器配置
	config := http.DefaultConfig()
	config.Addr = ":" + *port

	// 创建并启动HTTP服务器
	server := http.NewServer(config, gameState, staticFiles)
	if err := server.Start(); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
