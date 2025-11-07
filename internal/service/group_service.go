package service

import (
	"errors"

	"gorm.io/gorm"

	"miniblog/internal/model"
	"miniblog/internal/repository"
	"miniblog/internal/types"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
	userRepo  *repository.UserRepository
}

func NewGroupService(groupRepo *repository.GroupRepository, userRepo *repository.UserRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// CreateGroup 创建群组
func (s *GroupService) CreateGroup(ownerID int64, req *types.CreateGroupRequest) (*model.Group, error) {
	// 创建群组
	group := &model.Group{
		Name:        req.Name,
		Avatar:      req.Avatar,
		OwnerID:     ownerID,
		Description: req.Description,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// 添加群主
	ownerMember := &model.GroupMember{
		GroupID: group.ID,
		UserID:  ownerID,
		Role:    model.GroupRoleOwner,
	}
	if err := s.groupRepo.AddMember(ownerMember); err != nil {
		return nil, err
	}

	// 添加群成员
	for _, memberID := range req.MemberIDs {
		if memberID == ownerID {
			continue
		}
		member := &model.GroupMember{
			GroupID: group.ID,
			UserID:  memberID,
			Role:    model.GroupRoleMember,
		}
		if err := s.groupRepo.AddMember(member); err != nil {
			// 继续添加其他成员，不因为一个失败而终止
			continue
		}
	}

	return group, nil
}

// GetGroup 获取群组信息
func (s *GroupService) GetGroup(groupID int64) (*types.GroupResponse, error) {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}

	// 获取成员数量
	count, err := s.groupRepo.GetMemberCount(groupID)
	if err != nil {
		return nil, err
	}

	return &types.GroupResponse{
		Group:       *group,
		MemberCount: int(count),
	}, nil
}

// UpdateGroup 更新群组信息
func (s *GroupService) UpdateGroup(userID, groupID int64, req *types.UpdateGroupRequest) (*model.Group, error) {
	// 检查权限（只有群主和管理员可以修改）
	member, err := s.groupRepo.GetMember(groupID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not group member")
		}
		return nil, err
	}

	if member.Role != model.GroupRoleOwner && member.Role != model.GroupRoleAdmin {
		return nil, errors.New("no permission to update group")
	}

	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Avatar != "" {
		group.Avatar = req.Avatar
	}
	if req.Description != "" {
		group.Description = req.Description
	}

	if err := s.groupRepo.Update(group); err != nil {
		return nil, err
	}

	return group, nil
}

// DeleteGroup 删除群组
func (s *GroupService) DeleteGroup(userID, groupID int64) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("group not found")
		}
		return err
	}

	// 只有群主可以删除群组
	if group.OwnerID != userID {
		return errors.New("only group owner can delete group")
	}

	return s.groupRepo.Delete(groupID)
}

// GetUserGroups 获取用户加入的群组列表
func (s *GroupService) GetUserGroups(userID int64) ([]*types.GroupResponse, error) {
	groups, err := s.groupRepo.GetUserGroups(userID)
	if err != nil {
		return nil, err
	}

	var result []*types.GroupResponse
	for _, group := range groups {
		count, _ := s.groupRepo.GetMemberCount(group.ID)
		result = append(result, &types.GroupResponse{
			Group:       *group,
			MemberCount: int(count),
		})
	}

	return result, nil
}

// AddGroupMembers 添加群成员
func (s *GroupService) AddGroupMembers(userID, groupID int64, req *types.AddGroupMemberRequest) error {
	// 检查权限
	member, err := s.groupRepo.GetMember(groupID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not group member")
		}
		return err
	}

	if member.Role != model.GroupRoleOwner && member.Role != model.GroupRoleAdmin {
		return errors.New("no permission to add members")
	}

	// 添加成员
	for _, memberID := range req.UserIDs {
		// 检查是否已经是成员
		isMember, _ := s.groupRepo.IsMember(groupID, memberID)
		if isMember {
			continue
		}

		newMember := &model.GroupMember{
			GroupID: groupID,
			UserID:  memberID,
			Role:    model.GroupRoleMember,
		}
		if err := s.groupRepo.AddMember(newMember); err != nil {
			continue
		}
	}

	return nil
}

// RemoveGroupMember 移除群成员
func (s *GroupService) RemoveGroupMember(userID, groupID, memberID int64) error {
	// 检查权限
	operator, err := s.groupRepo.GetMember(groupID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not group member")
		}
		return err
	}

	if operator.Role != model.GroupRoleOwner && operator.Role != model.GroupRoleAdmin {
		return errors.New("no permission to remove members")
	}

	// 不能移除群主
	member, err := s.groupRepo.GetMember(groupID, memberID)
	if err != nil {
		return err
	}

	if member.Role == model.GroupRoleOwner {
		return errors.New("cannot remove group owner")
	}

	return s.groupRepo.RemoveMember(groupID, memberID)
}

// GetGroupMembers 获取群成员列表
func (s *GroupService) GetGroupMembers(groupID int64) ([]*types.GroupMemberResponse, error) {
	members, err := s.groupRepo.GetMembers(groupID)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	var userIDs []int64
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var result []*types.GroupMemberResponse
	for _, member := range members {
		if user, ok := userMap[member.UserID]; ok {
			result = append(result, &types.GroupMemberResponse{
				GroupMember: *member,
				UserInfo:    *user,
			})
		}
	}

	return result, nil
}

// LeaveGroup 退出群组
func (s *GroupService) LeaveGroup(userID, groupID int64) error {
	member, err := s.groupRepo.GetMember(groupID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not group member")
		}
		return err
	}

	// 群主不能直接退出，需要先转让群主
	if member.Role == model.GroupRoleOwner {
		return errors.New("group owner cannot leave, please transfer ownership first")
	}

	return s.groupRepo.RemoveMember(groupID, userID)
}

