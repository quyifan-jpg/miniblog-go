package types

import "miniblog/internal/model"

// LoginResponse 登录响应
type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// UserResponse 用户响应
type UserResponse struct {
	model.User
}

// FriendResponse 好友响应
type FriendResponse struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	FriendID  int64      `json:"friend_id"`
	Remark    string     `json:"remark"`
	Status    int8       `json:"status"`
	FriendInfo model.User `json:"friend_info"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	model.Message
	FromUserInfo *model.User `json:"from_user_info,omitempty"`
	ToUserInfo   *model.User `json:"to_user_info,omitempty"`
}

// GroupResponse 群组响应
type GroupResponse struct {
	model.Group
	MemberCount int `json:"member_count"`
}

// GroupMemberResponse 群成员响应
type GroupMemberResponse struct {
	model.GroupMember
	UserInfo model.User `json:"user_info"`
}

