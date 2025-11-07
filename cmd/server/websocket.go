package main

import (
	"miniblog/config"
	"miniblog/internal/pkg/jwt"
	"miniblog/internal/websocket"
)

// startWebSocketServer 启动 WebSocket 服务器
func startWebSocketServer(cfg *config.Config, jwtManager *jwt.JWTManager) error {
	server := websocket.NewServer(&cfg.Server.WebSocket, jwtManager)
	return server.Start()
}

