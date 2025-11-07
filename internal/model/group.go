package model

import (
	"time"
)

// Group 群组表
type Group struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Avatar      string    `gorm:"type:varchar(255)" json:"avatar"`
	OwnerID     int64     `gorm:"type:bigint;not null;index" json:"owner_id"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Group) TableName() string {
	return "groups"
}

// GroupMember 群成员表
type GroupMember struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID   int64     `gorm:"type:bigint;not null;index:idx_group_user,priority:1" json:"group_id"`
	UserID    int64     `gorm:"type:bigint;not null;index:idx_group_user,priority:2" json:"user_id"`
	Role      int8      `gorm:"type:tinyint;default:1" json:"role"` // 1-普通成员 2-管理员 3-群主
	JoinedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"joined_at"`
}

func (GroupMember) TableName() string {
	return "group_members"
}

const (
	GroupRoleMember = 1 // 普通成员
	GroupRoleAdmin  = 2 // 管理员
	GroupRoleOwner  = 3 // 群主
)

