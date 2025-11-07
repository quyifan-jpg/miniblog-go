package model

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex:idx_username;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // 不返回密码
	Nickname  string    `gorm:"type:varchar(100)" json:"nickname"`
	Avatar    string    `gorm:"type:varchar(255)" json:"avatar"`
	Email     string    `gorm:"type:varchar(100)" json:"email"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

