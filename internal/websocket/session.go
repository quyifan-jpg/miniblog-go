package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"miniblog/internal/pkg/logger"
	"miniblog/internal/types"
)

// Session WebSocket 会话
type Session struct {
	ID          string
	UserID      int64
	Conn        *websocket.Conn
	Send        chan []byte
	manager     *Manager
	lastPing    time.Time
	mu          sync.RWMutex
	isClosed    bool
}

// NewSession 创建新会话
func NewSession(id string, userID int64, conn *websocket.Conn, manager *Manager) *Session {
	return &Session{
		ID:       id,
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		manager:  manager,
		lastPing: time.Now(),
	}
}

// ReadPump 读取消息
func (s *Session) ReadPump() {
	defer func() {
		s.manager.Unregister(s)
		s.Close()
	}()

	s.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	s.Conn.SetPongHandler(func(string) error {
		s.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		s.UpdateLastPing()
		return nil
	})

	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket read error", zap.Error(err))
			}
			break
		}

		// 处理接收到的消息
		s.handleMessage(message)
	}
}

// WritePump 发送消息
func (s *Session) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		s.Close()
	}()

	for {
		select {
		case message, ok := <-s.Send:
			s.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				s.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := s.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Error("WebSocket write error", zap.Error(err))
				return
			}

		case <-ticker.C:
			s.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理消息
func (s *Session) handleMessage(message []byte) {
	var wsMsg types.WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		logger.Error("Failed to unmarshal message", zap.Error(err))
		return
	}

	switch wsMsg.Event {
	case types.EventTypePing:
		// 回复 pong
		s.SendMessage(types.EventTypePong, nil)
	case types.EventTypeMessage:
		// 消息已通过 HTTP API 发送，这里只是确认
		logger.Info("Received message event", zap.Any("data", wsMsg.Data))
	default:
		logger.Warn("Unknown event type", zap.String("event", wsMsg.Event))
	}
}

// SendMessage 发送消息
func (s *Session) SendMessage(event string, data interface{}) error {
	wsMsg := types.WSMessage{
		Event:     event,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	msgBytes, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	select {
	case s.Send <- msgBytes:
		return nil
	default:
		logger.Warn("Session send channel is full", zap.Int64("user_id", s.UserID))
		return nil
	}
}

// UpdateLastPing 更新最后 ping 时间
func (s *Session) UpdateLastPing() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastPing = time.Now()
}

// GetLastPing 获取最后 ping 时间
func (s *Session) GetLastPing() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastPing
}

// Close 关闭会话
func (s *Session) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isClosed {
		return
	}

	s.isClosed = true
	close(s.Send)
	s.Conn.Close()
}

// IsClosed 是否已关闭
func (s *Session) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isClosed
}

