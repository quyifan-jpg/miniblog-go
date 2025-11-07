package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"miniblog/internal/middleware"
	"miniblog/internal/pkg/response"
	"miniblog/internal/service"
	"miniblog/internal/types"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

// AddFriend 添加好友
// @Summary 添加好友
// @Tags friend
// @Accept json
// @Produce json
// @Param request body types.AddFriendRequest true "好友信息"
// @Success 200 {object} response.Response
// @Router /api/friend/add [post]
func (h *FriendHandler) AddFriend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req types.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	if err := h.friendService.AddFriend(userID, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// HandleFriendRequest 处理好友请求
// @Summary 处理好友请求
// @Tags friend
// @Accept json
// @Produce json
// @Param request body types.HandleFriendRequest true "处理信息"
// @Success 200 {object} response.Response
// @Router /api/friend/handle [post]
func (h *FriendHandler) HandleFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req types.HandleFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	if err := h.friendService.HandleFriendRequest(userID, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetFriendList 获取好友列表
// @Summary 获取好友列表
// @Tags friend
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/friend/list [get]
func (h *FriendHandler) GetFriendList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	friends, err := h.friendService.GetFriendList(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, friends)
}

// GetFriendRequests 获取好友请求列表
// @Summary 获取好友请求列表
// @Tags friend
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/friend/requests [get]
func (h *FriendHandler) GetFriendRequests(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	requests, err := h.friendService.GetFriendRequests(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, requests)
}

// DeleteFriend 删除好友
// @Summary 删除好友
// @Tags friend
// @Produce json
// @Param id path int true "好友ID"
// @Success 200 {object} response.Response
// @Router /api/friend/:id [delete]
func (h *FriendHandler) DeleteFriend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	friendIDStr := c.Param("id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid friend id")
		return
	}

	if err := h.friendService.DeleteFriend(userID, friendID); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

