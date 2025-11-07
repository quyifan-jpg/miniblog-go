package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"miniblog/internal/middleware"
	"miniblog/internal/pkg/response"
	"miniblog/internal/service"
	"miniblog/internal/types"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup 创建群组
// @Summary 创建群组
// @Tags group
// @Accept json
// @Produce json
// @Param request body types.CreateGroupRequest true "群组信息"
// @Success 200 {object} response.Response
// @Router /api/group/create [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req types.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	group, err := h.groupService.CreateGroup(userID, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, group)
}

// GetGroup 获取群组信息
// @Summary 获取群组信息
// @Tags group
// @Produce json
// @Param id path int true "群组ID"
// @Success 200 {object} response.Response
// @Router /api/group/:id [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	group, err := h.groupService.GetGroup(groupID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, group)
}

// UpdateGroup 更新群组信息
// @Summary 更新群组信息
// @Tags group
// @Accept json
// @Produce json
// @Param id path int true "群组ID"
// @Param request body types.UpdateGroupRequest true "群组信息"
// @Success 200 {object} response.Response
// @Router /api/group/:id [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	var req types.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	group, err := h.groupService.UpdateGroup(userID, groupID, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, group)
}

// DeleteGroup 删除群组
// @Summary 删除群组
// @Tags group
// @Produce json
// @Param id path int true "群组ID"
// @Success 200 {object} response.Response
// @Router /api/group/:id [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	if err := h.groupService.DeleteGroup(userID, groupID); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUserGroups 获取用户加入的群组列表
// @Summary 获取用户加入的群组列表
// @Tags group
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/group/list [get]
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groups, err := h.groupService.GetUserGroups(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, groups)
}

// AddGroupMembers 添加群成员
// @Summary 添加群成员
// @Tags group
// @Accept json
// @Produce json
// @Param id path int true "群组ID"
// @Param request body types.AddGroupMemberRequest true "成员信息"
// @Success 200 {object} response.Response
// @Router /api/group/:id/member/add [post]
func (h *GroupHandler) AddGroupMembers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	var req types.AddGroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	if err := h.groupService.AddGroupMembers(userID, groupID, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// RemoveGroupMember 移除群成员
// @Summary 移除群成员
// @Tags group
// @Produce json
// @Param id path int true "群组ID"
// @Param member_id path int true "成员ID"
// @Success 200 {object} response.Response
// @Router /api/group/:id/member/:member_id [delete]
func (h *GroupHandler) RemoveGroupMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	memberIDStr := c.Param("member_id")
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid member id")
		return
	}

	if err := h.groupService.RemoveGroupMember(userID, groupID, memberID); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetGroupMembers 获取群成员列表
// @Summary 获取群成员列表
// @Tags group
// @Produce json
// @Param id path int true "群组ID"
// @Success 200 {object} response.Response
// @Router /api/group/:id/members [get]
func (h *GroupHandler) GetGroupMembers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	members, err := h.groupService.GetGroupMembers(groupID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, members)
}

// LeaveGroup 退出群组
// @Summary 退出群组
// @Tags group
// @Produce json
// @Param id path int true "群组ID"
// @Success 200 {object} response.Response
// @Router /api/group/:id/leave [post]
func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	if err := h.groupService.LeaveGroup(userID, groupID); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

