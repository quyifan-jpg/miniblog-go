package types

// WebSocket 事件类型
const (
	EventTypePing         = "ping"
	EventTypePong         = "pong"
	EventTypeMessage      = "message"
	EventTypeMessageRead  = "message_read"
	EventTypeFriendOnline = "friend_online"
	EventTypeFriendOffline = "friend_offline"
	EventTypeGroupMessage = "group_message"
)

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
	MsgID   string      `json:"msg_id,omitempty"`
	Timestamp int64     `json:"timestamp"`
}

// MessageEventData 消息事件数据
type MessageEventData struct {
	ID         int64  `json:"id"`
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
	GroupID    int64  `json:"group_id"`
	MsgType    int8   `json:"msg_type"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

// OnlineStatusData 在线状态数据
type OnlineStatusData struct {
	UserID int64 `json:"user_id"`
	Online bool  `json:"online"`
}

