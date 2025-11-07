package repository

import (
	"miniblog/internal/model"
	"miniblog/internal/pkg/mysql"
)

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

// Create 创建消息
func (r *MessageRepository) Create(message *model.Message) error {
	return mysql.DB.Create(message).Error
}

// GetByID 根据 ID 获取消息
func (r *MessageRepository) GetByID(id int64) (*model.Message, error) {
	var message model.Message
	err := mysql.DB.Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// Update 更新消息
func (r *MessageRepository) Update(message *model.Message) error {
	return mysql.DB.Save(message).Error
}

// Delete 删除消息
func (r *MessageRepository) Delete(id int64) error {
	return mysql.DB.Delete(&model.Message{}, id).Error
}

// GetChatHistory 获取两个用户之间的聊天历史
func (r *MessageRepository) GetChatHistory(userID1, userID2 int64, offset, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	err := mysql.DB.Where(
		"(from_user_id = ? AND to_user_id = ? AND group_id = 0) OR (from_user_id = ? AND to_user_id = ? AND group_id = 0)",
		userID1, userID2, userID2, userID1,
	).Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error
	return messages, err
}

// GetGroupChatHistory 获取群聊历史
func (r *MessageRepository) GetGroupChatHistory(groupID int64, offset, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	err := mysql.DB.Where("group_id = ?", groupID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error
	return messages, err
}

// GetUnreadMessages 获取未读消息
func (r *MessageRepository) GetUnreadMessages(userID int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := mysql.DB.Where("to_user_id = ? AND status = ?", userID, model.MsgStatusUnread).
		Order("created_at DESC").Find(&messages).Error
	return messages, err
}

// MarkAsRead 标记消息为已读
func (r *MessageRepository) MarkAsRead(messageIDs []int64) error {
	return mysql.DB.Model(&model.Message{}).Where("id IN ?", messageIDs).
		Update("status", model.MsgStatusRead).Error
}

// MarkChatAsRead 标记聊天消息为已读
func (r *MessageRepository) MarkChatAsRead(fromUserID, toUserID int64) error {
	return mysql.DB.Model(&model.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND status = ?", fromUserID, toUserID, model.MsgStatusUnread).
		Update("status", model.MsgStatusRead).Error
}

// GetUnreadCount 获取未读消息数量
func (r *MessageRepository) GetUnreadCount(userID int64) (int64, error) {
	var count int64
	err := mysql.DB.Model(&model.Message{}).
		Where("to_user_id = ? AND status = ?", userID, model.MsgStatusUnread).
		Count(&count).Error
	return count, err
}

