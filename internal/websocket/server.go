package websocket

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"miniblog/config"
	"miniblog/internal/pkg/jwt"
	"miniblog/internal/pkg/logger"
)

// Server WebSocket 服务器
type Server struct {
	manager *Manager
	handler *Handler
	cfg     *config.WebSocketConfig
}

// NewServer 创建 WebSocket 服务器
func NewServer(cfg *config.WebSocketConfig, jwtManager *jwt.JWTManager) *Server {
	manager := NewManager()
	handler := NewHandler(manager, jwtManager)

	return &Server{
		manager: manager,
		handler: handler,
		cfg:     cfg,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 启动管理器
	go s.manager.Run()

	// 创建 Gin 路由
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// WebSocket 路由
	router.GET("/ws", s.handler.HandleWebSocket)

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"online_users": len(s.manager.GetOnlineUserIDs()),
		})
	})

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	logger.Info("WebSocket server starting", zap.String("addr", addr))

	return http.ListenAndServe(addr, router)
}

// Stop 停止服务器
func (s *Server) Stop() {
	s.manager.Stop()
	logger.Info("WebSocket server stopped")
}

// GetManager 获取管理器
func (s *Server) GetManager() *Manager {
	return s.manager
}

