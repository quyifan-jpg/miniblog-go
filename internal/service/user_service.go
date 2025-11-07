package service

import (
	"errors"

	"gorm.io/gorm"

	"miniblog/internal/model"
	"miniblog/internal/repository"
	"miniblog/internal/types"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID int64, req *types.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// SearchUsers 搜索用户
func (s *UserService) SearchUsers(keyword string, offset, limit int) ([]*model.User, error) {
	return s.userRepo.Search(keyword, offset, limit)
}

// GetUsersByIDs 根据 ID 列表获取用户
func (s *UserService) GetUsersByIDs(ids []int64) ([]*model.User, error) {
	return s.userRepo.GetByIDs(ids)
}

