package types

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"max=100"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// AddFriendRequest 添加好友请求
type AddFriendRequest struct {
	FriendID int64  `json:"friend_id" binding:"required,gt=0"`
	Remark   string `json:"remark" binding:"max=100"`
}

// HandleFriendRequest 处理好友请求
type HandleFriendRequest struct {
	FriendID int64 `json:"friend_id" binding:"required,gt=0"`
	Accept   bool  `json:"accept"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ToUserID int64  `json:"to_user_id" binding:"required,gt=0"`
	GroupID  int64  `json:"group_id"`
	MsgType  int8   `json:"msg_type" binding:"required,oneof=1 2 3 4"`
	Content  string `json:"content" binding:"required"`
}

// CreateGroupRequest 创建群组请求
type CreateGroupRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Avatar      string  `json:"avatar" binding:"max=255"`
	Description string  `json:"description" binding:"max=500"`
	MemberIDs   []int64 `json:"member_ids"`
}

// UpdateGroupRequest 更新群组请求
type UpdateGroupRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Avatar      string `json:"avatar" binding:"max=255"`
	Description string `json:"description" binding:"max=500"`
}

// AddGroupMemberRequest 添加群成员请求
type AddGroupMemberRequest struct {
	UserIDs []int64 `json:"user_ids" binding:"required,min=1"`
}

