package service

import (
	"errors"

	"gorm.io/gorm"

	"miniblog/internal/model"
	"miniblog/internal/repository"
	"miniblog/internal/types"
)

type FriendService struct {
	friendRepo *repository.FriendRepository
	userRepo   *repository.UserRepository
}

func NewFriendService(friendRepo *repository.FriendRepository, userRepo *repository.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

// AddFriend 添加好友请求
func (s *FriendService) AddFriend(userID int64, req *types.AddFriendRequest) error {
	// 检查好友是否存在
	_, err := s.userRepo.GetByID(req.FriendID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friend not found")
		}
		return err
	}

	// 不能添加自己为好友
	if userID == req.FriendID {
		return errors.New("cannot add yourself as friend")
	}

	// 检查是否已经是好友或有待处理的请求
	existFriend, err := s.friendRepo.GetByUserAndFriend(userID, req.FriendID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existFriend != nil {
		if existFriend.Status == model.FriendStatusAccepted {
			return errors.New("already friends")
		}
		if existFriend.Status == model.FriendStatusPending {
			return errors.New("friend request already sent")
		}
	}

	// 创建好友请求
	friend := &model.Friend{
		UserID:   userID,
		FriendID: req.FriendID,
		Remark:   req.Remark,
		Status:   model.FriendStatusPending,
	}

	return s.friendRepo.Create(friend)
}

// HandleFriendRequest 处理好友请求
func (s *FriendService) HandleFriendRequest(userID int64, req *types.HandleFriendRequest) error {
	// 获取好友请求
	friend, err := s.friendRepo.GetByUserAndFriend(req.FriendID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friend request not found")
		}
		return err
	}

	if friend.Status != model.FriendStatusPending {
		return errors.New("friend request already handled")
	}

	// 更新状态
	if req.Accept {
		friend.Status = model.FriendStatusAccepted

		// 创建双向好友关系
		reverseFriend := &model.Friend{
			UserID:   userID,
			FriendID: req.FriendID,
			Status:   model.FriendStatusAccepted,
		}
		if err := s.friendRepo.Create(reverseFriend); err != nil {
			return err
		}
	} else {
		friend.Status = model.FriendStatusRejected
	}

	return s.friendRepo.Update(friend)
}

// GetFriendList 获取好友列表
func (s *FriendService) GetFriendList(userID int64) ([]*types.FriendResponse, error) {
	friends, err := s.friendRepo.GetFriendList(userID, model.FriendStatusAccepted)
	if err != nil {
		return nil, err
	}

	// 获取好友信息
	var friendIDs []int64
	for _, friend := range friends {
		friendIDs = append(friendIDs, friend.FriendID)
	}

	users, err := s.userRepo.GetByIDs(friendIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var result []*types.FriendResponse
	for _, friend := range friends {
		if user, ok := userMap[friend.FriendID]; ok {
			result = append(result, &types.FriendResponse{
				ID:         friend.ID,
				UserID:     friend.UserID,
				FriendID:   friend.FriendID,
				Remark:     friend.Remark,
				Status:     friend.Status,
				FriendInfo: *user,
			})
		}
	}

	return result, nil
}

// GetFriendRequests 获取好友请求列表
func (s *FriendService) GetFriendRequests(userID int64) ([]*types.FriendResponse, error) {
	friends, err := s.friendRepo.GetFriendRequests(userID)
	if err != nil {
		return nil, err
	}

	// 获取请求者信息
	var userIDs []int64
	for _, friend := range friends {
		userIDs = append(userIDs, friend.UserID)
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var result []*types.FriendResponse
	for _, friend := range friends {
		if user, ok := userMap[friend.UserID]; ok {
			result = append(result, &types.FriendResponse{
				ID:         friend.ID,
				UserID:     friend.UserID,
				FriendID:   friend.FriendID,
				Remark:     friend.Remark,
				Status:     friend.Status,
				FriendInfo: *user,
			})
		}
	}

	return result, nil
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(userID, friendID int64) error {
	// 删除双向好友关系
	if err := s.friendRepo.DeleteByUserAndFriend(userID, friendID); err != nil {
		return err
	}
	if err := s.friendRepo.DeleteByUserAndFriend(friendID, userID); err != nil {
		return err
	}
	return nil
}

// IsFriend 判断是否是好友
func (s *FriendService) IsFriend(userID, friendID int64) (bool, error) {
	return s.friendRepo.IsFriend(userID, friendID)
}

