package websocket

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"miniblog/internal/pkg/jwt"
	"miniblog/internal/pkg/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该做严格检查
	},
}

// Handler WebSocket 处理器
type Handler struct {
	manager    *Manager
	jwtManager *jwt.JWTManager
}

// NewHandler 创建处理器
func NewHandler(manager *Manager, jwtManager *jwt.JWTManager) *Handler {
	return &Handler{
		manager:    manager,
		jwtManager: jwtManager,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *Handler) HandleWebSocket(c *gin.Context) {
	// 从查询参数或 Header 中获取 Token
	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	// 验证 Token
	claims, err := h.jwtManager.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade to websocket", zap.Error(err))
		return
	}

	// 创建会话
	sessionID := uuid.New().String()
	session := NewSession(sessionID, claims.UserID, conn, h.manager)

	// 注册会话
	h.manager.Register(session)

	// 启动读写协程
	go session.WritePump()
	go session.ReadPump()

	logger.Info("WebSocket connection established", 
		zap.String("session_id", sessionID),
		zap.Int64("user_id", claims.UserID),
	)
}

