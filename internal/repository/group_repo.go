package repository

import (
	"miniblog/internal/model"
	"miniblog/internal/pkg/mysql"
)

type GroupRepository struct{}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{}
}

// Create 创建群组
func (r *GroupRepository) Create(group *model.Group) error {
	return mysql.DB.Create(group).Error
}

// GetByID 根据 ID 获取群组
func (r *GroupRepository) GetByID(id int64) (*model.Group, error) {
	var group model.Group
	err := mysql.DB.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// Update 更新群组信息
func (r *GroupRepository) Update(group *model.Group) error {
	return mysql.DB.Save(group).Error
}

// Delete 删除群组
func (r *GroupRepository) Delete(id int64) error {
	return mysql.DB.Delete(&model.Group{}, id).Error
}

// GetUserGroups 获取用户加入的群组列表
func (r *GroupRepository) GetUserGroups(userID int64) ([]*model.Group, error) {
	var groups []*model.Group
	err := mysql.DB.Table("groups").
		Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

// AddMember 添加群成员
func (r *GroupRepository) AddMember(member *model.GroupMember) error {
	return mysql.DB.Create(member).Error
}

// RemoveMember 移除群成员
func (r *GroupRepository) RemoveMember(groupID, userID int64) error {
	return mysql.DB.Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.GroupMember{}).Error
}

// GetMember 获取群成员
func (r *GroupRepository) GetMember(groupID, userID int64) (*model.GroupMember, error) {
	var member model.GroupMember
	err := mysql.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMembers 获取群成员列表
func (r *GroupRepository) GetMembers(groupID int64) ([]*model.GroupMember, error) {
	var members []*model.GroupMember
	err := mysql.DB.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

// GetMemberIDs 获取群成员 ID 列表
func (r *GroupRepository) GetMemberIDs(groupID int64) ([]int64, error) {
	var memberIDs []int64
	err := mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &memberIDs).Error
	return memberIDs, err
}

// IsMember 判断是否是群成员
func (r *GroupRepository) IsMember(groupID, userID int64) (bool, error) {
	var count int64
	err := mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

// UpdateMemberRole 更新群成员角色
func (r *GroupRepository) UpdateMemberRole(groupID, userID int64, role int8) error {
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", role).Error
}

// GetMemberCount 获取群成员数量
func (r *GroupRepository) GetMemberCount(groupID int64) (int64, error) {
	var count int64
	err := mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return count, err
}

