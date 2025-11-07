package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"

	"miniblog/internal/pkg/logger"
	"miniblog/internal/pkg/redis"
	"miniblog/internal/types"
)

// Manager WebSocket 连接管理器
type Manager struct {
	sessions   map[int64]*Session // userID -> Session
	register   chan *Session
	unregister chan *Session
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewManager 创建管理器
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		sessions:   make(map[int64]*Session),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Run 运行管理器
func (m *Manager) Run() {
	// 启动心跳检测
	go m.heartbeatCheck()

	// 订阅 Redis 消息
	go m.subscribeMessages()

	for {
		select {
		case session := <-m.register:
			m.addSession(session)
		case session := <-m.unregister:
			m.removeSession(session)
		case <-m.ctx.Done():
			return
		}
	}
}

// Register 注册会话
func (m *Manager) Register(session *Session) {
	m.register <- session
}

// Unregister 注销会话
func (m *Manager) Unregister(session *Session) {
	m.unregister <- session
}

// addSession 添加会话
func (m *Manager) addSession(session *Session) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果用户已有连接，关闭旧连接
	if oldSession, ok := m.sessions[session.UserID]; ok {
		oldSession.Close()
	}

	m.sessions[session.UserID] = session
	logger.Info("Session registered", zap.Int64("user_id", session.UserID))

	// 通知好友上线
	m.notifyFriendsOnline(session.UserID, true)
}

// removeSession 移除会话
func (m *Manager) removeSession(session *Session) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[session.UserID]; ok {
		delete(m.sessions, session.UserID)
		logger.Info("Session unregistered", zap.Int64("user_id", session.UserID))

		// 通知好友离线
		m.notifyFriendsOnline(session.UserID, false)
	}
}

// GetSession 获取会话
func (m *Manager) GetSession(userID int64) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[userID]
	return session, ok
}

// SendToUser 发送消息给用户
func (m *Manager) SendToUser(userID int64, event string, data interface{}) error {
	session, ok := m.GetSession(userID)
	if !ok {
		logger.Warn("User session not found", zap.Int64("user_id", userID))
		return nil
	}

	return session.SendMessage(event, data)
}

// SendToUsers 发送消息给多个用户
func (m *Manager) SendToUsers(userIDs []int64, event string, data interface{}) {
	for _, userID := range userIDs {
		m.SendToUser(userID, event, data)
	}
}

// Broadcast 广播消息
func (m *Manager) Broadcast(event string, data interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, session := range m.sessions {
		session.SendMessage(event, data)
	}
}

// IsOnline 检查用户是否在线
func (m *Manager) IsOnline(userID int64) bool {
	_, ok := m.GetSession(userID)
	return ok
}

// GetOnlineUserIDs 获取所有在线用户 ID
func (m *Manager) GetOnlineUserIDs() []int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userIDs := make([]int64, 0, len(m.sessions))
	for userID := range m.sessions {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// heartbeatCheck 心跳检测
func (m *Manager) heartbeatCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkDeadConnections()
		case <-m.ctx.Done():
			return
		}
	}
}

// checkDeadConnections 检查死连接
func (m *Manager) checkDeadConnections() {
	m.mu.RLock()
	var deadSessions []*Session
	now := time.Now()

	for _, session := range m.sessions {
		if now.Sub(session.GetLastPing()) > 90*time.Second {
			deadSessions = append(deadSessions, session)
		}
	}
	m.mu.RUnlock()

	for _, session := range deadSessions {
		logger.Warn("Dead connection detected", zap.Int64("user_id", session.UserID))
		m.Unregister(session)
	}
}

// subscribeMessages 订阅 Redis 消息
func (m *Manager) subscribeMessages() {
	pubsub := redis.Subscribe(m.ctx, "im:message", "im:group:message")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			m.handleRedisMessage([]byte(msg.Payload))
		case <-m.ctx.Done():
			return
		}
	}
}

// handleRedisMessage 处理 Redis 消息
func (m *Manager) handleRedisMessage(data []byte) {
	var wsMsg types.WSMessage
	if err := json.Unmarshal(data, &wsMsg); err != nil {
		logger.Error("Failed to unmarshal redis message", zap.Error(err))
		return
	}

	switch wsMsg.Event {
	case types.EventTypeMessage:
		// 单聊消息
		if msgData, ok := wsMsg.Data.(map[string]interface{}); ok {
			if toUserID, ok := msgData["to_user_id"].(float64); ok {
				m.SendToUser(int64(toUserID), wsMsg.Event, wsMsg.Data)
			}
		}
	case types.EventTypeGroupMessage:
		// 群聊消息 - 需要获取群成员并发送
		// 这里简化处理，实际应该从数据库获取群成员列表
		logger.Info("Group message received", zap.Any("data", wsMsg.Data))
	}
}

// notifyFriendsOnline 通知好友上线/离线
func (m *Manager) notifyFriendsOnline(userID int64, online bool) {
	// 这里应该从数据库获取好友列表，简化处理
	data := types.OnlineStatusData{
		UserID: userID,
		Online: online,
	}

	event := types.EventTypeFriendOnline
	if !online {
		event = types.EventTypeFriendOffline
	}

	// 发布到 Redis，由其他服务实例处理
	wsMsg := types.WSMessage{
		Event:     event,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	msgBytes, _ := json.Marshal(wsMsg)
	redis.Publish(m.ctx, "im:online:status", msgBytes)
}

// Stop 停止管理器
func (m *Manager) Stop() {
	m.cancel()

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, session := range m.sessions {
		session.Close()
	}
}

