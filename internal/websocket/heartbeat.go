package websocket

import (
	"time"
)

// HeartbeatConfig 心跳配置
type HeartbeatConfig struct {
	PingInterval  time.Duration // Ping 间隔
	PongTimeout   time.Duration // Pong 超时时间
	CheckInterval time.Duration // 检查间隔
}

// DefaultHeartbeatConfig 默认心跳配置
var DefaultHeartbeatConfig = HeartbeatConfig{
	PingInterval:  30 * time.Second,
	PongTimeout:   60 * time.Second,
	CheckInterval: 30 * time.Second,
}

// IsAlive 检查会话是否存活
func (s *Session) IsAlive(timeout time.Duration) bool {
	return time.Since(s.GetLastPing()) <= timeout
}

