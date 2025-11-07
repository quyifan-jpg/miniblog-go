package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"miniblog/internal/model"
	"miniblog/internal/pkg/redis"
	"miniblog/internal/repository"
	"miniblog/internal/types"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
	friendRepo  *repository.FriendRepository
	groupRepo   *repository.GroupRepository
}

func NewMessageService(
	messageRepo *repository.MessageRepository,
	userRepo *repository.UserRepository,
	friendRepo *repository.FriendRepository,
	groupRepo *repository.GroupRepository,
) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		friendRepo:  friendRepo,
		groupRepo:   groupRepo,
	}
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(fromUserID int64, req *types.SendMessageRequest) (*model.Message, error) {
	// 单聊消息
	if req.GroupID == 0 {
		// 检查是否是好友
		isFriend, err := s.friendRepo.IsFriend(fromUserID, req.ToUserID)
		if err != nil {
			return nil, err
		}
		if !isFriend {
			return nil, errors.New("not friends, cannot send message")
		}
	} else {
		// 群聊消息，检查是否是群成员
		isMember, err := s.groupRepo.IsMember(req.GroupID, fromUserID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("not group member, cannot send message")
		}
	}

	// 创建消息
	message := &model.Message{
		FromUserID: fromUserID,
		ToUserID:   req.ToUserID,
		GroupID:    req.GroupID,
		MsgType:    req.MsgType,
		Content:    req.Content,
		Status:     model.MsgStatusUnread,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	// 发布消息到 Redis，供 WebSocket 服务使用
	if err := s.publishMessage(message); err != nil {
		// 记录日志，但不影响消息发送
		// logger.Error("failed to publish message", zap.Error(err))
	}

	return message, nil
}

// GetChatHistory 获取聊天历史
func (s *MessageService) GetChatHistory(userID, friendID int64, offset, limit int) ([]*types.MessageResponse, error) {
	messages, err := s.messageRepo.GetChatHistory(userID, friendID, offset, limit)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

// GetGroupChatHistory 获取群聊历史
func (s *MessageService) GetGroupChatHistory(groupID int64, offset, limit int) ([]*types.MessageResponse, error) {
	messages, err := s.messageRepo.GetGroupChatHistory(groupID, offset, limit)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

// GetUnreadMessages 获取未读消息
func (s *MessageService) GetUnreadMessages(userID int64) ([]*types.MessageResponse, error) {
	messages, err := s.messageRepo.GetUnreadMessages(userID)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

// MarkAsRead 标记消息为已读
func (s *MessageService) MarkAsRead(messageIDs []int64) error {
	return s.messageRepo.MarkAsRead(messageIDs)
}

// MarkChatAsRead 标记聊天消息为已读
func (s *MessageService) MarkChatAsRead(fromUserID, toUserID int64) error {
	return s.messageRepo.MarkChatAsRead(fromUserID, toUserID)
}

// GetUnreadCount 获取未读消息数量
func (s *MessageService) GetUnreadCount(userID int64) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}

// buildMessageResponses 构建消息响应
func (s *MessageService) buildMessageResponses(messages []*model.Message) ([]*types.MessageResponse, error) {
	if len(messages) == 0 {
		return []*types.MessageResponse{}, nil
	}

	// 收集所有用户 ID
	userIDMap := make(map[int64]bool)
	for _, msg := range messages {
		userIDMap[msg.FromUserID] = true
		if msg.ToUserID > 0 {
			userIDMap[msg.ToUserID] = true
		}
	}

	var userIDs []int64
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	// 获取用户信息
	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 构建响应
	var result []*types.MessageResponse
	for _, msg := range messages {
		resp := &types.MessageResponse{
			Message: *msg,
		}
		if user, ok := userMap[msg.FromUserID]; ok {
			resp.FromUserInfo = user
		}
		if msg.ToUserID > 0 {
			if user, ok := userMap[msg.ToUserID]; ok {
				resp.ToUserInfo = user
			}
		}
		result = append(result, resp)
	}

	return result, nil
}

// publishMessage 发布消息到 Redis
func (s *MessageService) publishMessage(message *model.Message) error {
	eventData := types.MessageEventData{
		ID:         message.ID,
		FromUserID: message.FromUserID,
		ToUserID:   message.ToUserID,
		GroupID:    message.GroupID,
		MsgType:    message.MsgType,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt.Format(time.RFC3339),
	}

	wsMsg := types.WSMessage{
		Event:     types.EventTypeMessage,
		Data:      eventData,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	ctx := context.Background()
	channel := "im:message"
	if message.GroupID > 0 {
		channel = "im:group:message"
	}

	return redis.Publish(ctx, channel, data)
}

