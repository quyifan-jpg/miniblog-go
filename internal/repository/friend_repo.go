package repository

import (
	"miniblog/internal/model"
	"miniblog/internal/pkg/mysql"
)

type FriendRepository struct{}

func NewFriendRepository() *FriendRepository {
	return &FriendRepository{}
}

// Create 创建好友关系
func (r *FriendRepository) Create(friend *model.Friend) error {
	return mysql.DB.Create(friend).Error
}

// GetByID 根据 ID 获取好友关系
func (r *FriendRepository) GetByID(id int64) (*model.Friend, error) {
	var friend model.Friend
	err := mysql.DB.Where("id = ?", id).First(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

// GetByUserAndFriend 获取两个用户之间的好友关系
func (r *FriendRepository) GetByUserAndFriend(userID, friendID int64) (*model.Friend, error) {
	var friend model.Friend
	err := mysql.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

// Update 更新好友关系
func (r *FriendRepository) Update(friend *model.Friend) error {
	return mysql.DB.Save(friend).Error
}

// Delete 删除好友关系
func (r *FriendRepository) Delete(id int64) error {
	return mysql.DB.Delete(&model.Friend{}, id).Error
}

// DeleteByUserAndFriend 删除两个用户之间的好友关系
func (r *FriendRepository) DeleteByUserAndFriend(userID, friendID int64) error {
	return mysql.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&model.Friend{}).Error
}

// GetFriendList 获取用户的好友列表
func (r *FriendRepository) GetFriendList(userID int64, status int8) ([]*model.Friend, error) {
	var friends []*model.Friend
	query := mysql.DB.Where("user_id = ?", userID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&friends).Error
	return friends, err
}

// GetFriendRequests 获取好友请求列表
func (r *FriendRepository) GetFriendRequests(userID int64) ([]*model.Friend, error) {
	var friends []*model.Friend
	err := mysql.DB.Where("friend_id = ? AND status = ?", userID, model.FriendStatusPending).Find(&friends).Error
	return friends, err
}

// GetFriendIDs 获取用户的好友 ID 列表
func (r *FriendRepository) GetFriendIDs(userID int64) ([]int64, error) {
	var friendIDs []int64
	err := mysql.DB.Model(&model.Friend{}).
		Where("user_id = ? AND status = ?", userID, model.FriendStatusAccepted).
		Pluck("friend_id", &friendIDs).Error
	return friendIDs, err
}

// IsFriend 判断是否是好友
func (r *FriendRepository) IsFriend(userID, friendID int64) (bool, error) {
	var count int64
	err := mysql.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, model.FriendStatusAccepted).
		Count(&count).Error
	return count > 0, err
}

