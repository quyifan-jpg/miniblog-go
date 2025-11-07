package model

import (
	"time"
)

// Message 消息表
type Message struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserID int64     `gorm:"type:bigint;not null;index:idx_from" json:"from_user_id"`
	ToUserID   int64     `gorm:"type:bigint;not null;index:idx_to" json:"to_user_id"`
	GroupID    int64     `gorm:"type:bigint;default:0;index:idx_group" json:"group_id"` // 0表示单聊
	MsgType    int8      `gorm:"type:tinyint;not null" json:"msg_type"`                 // 1-文本 2-图片 3-语音 4-视频
	Content    string    `gorm:"type:text;not null" json:"content"`
	Status     int8      `gorm:"type:tinyint;default:0" json:"status"`  // 0-未读 1-已读
	CreatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

const (
	MsgTypeText  = 1 // 文本消息
	MsgTypeImage = 2 // 图片消息
	MsgTypeVoice = 3 // 语音消息
	MsgTypeVideo = 4 // 视频消息
)

const (
	MsgStatusUnread = 0 // 未读
	MsgStatusRead   = 1 // 已读
)

