package model

import (
	"time"
)

// Friend 好友关系表
type Friend struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"type:bigint;not null;index:idx_user_friend,priority:1" json:"user_id"`
	FriendID  int64     `gorm:"type:bigint;not null;index:idx_user_friend,priority:2" json:"friend_id"`
	Remark    string    `gorm:"type:varchar(100)" json:"remark"`       // 备注名
	Status    int8      `gorm:"type:tinyint;default:0" json:"status"`  // 0-待确认 1-已同意 2-已拒绝
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Friend) TableName() string {
	return "friends"
}

const (
	FriendStatusPending  = 0 // 待确认
	FriendStatusAccepted = 1 // 已同意
	FriendStatusRejected = 2 // 已拒绝
)

