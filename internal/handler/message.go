package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"miniblog/internal/middleware"
	"miniblog/internal/pkg/response"
	"miniblog/internal/service"
	"miniblog/internal/types"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// SendMessage 发送消息
// @Summary 发送消息
// @Tags message
// @Accept json
// @Produce json
// @Param request body types.SendMessageRequest true "消息信息"
// @Success 200 {object} response.Response
// @Router /api/message/send [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req types.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	message, err := h.messageService.SendMessage(userID, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, message)
}

// GetChatHistory 获取聊天历史
// @Summary 获取聊天历史
// @Tags message
// @Produce json
// @Param friend_id query int true "好友ID"
// @Param offset query int false "偏移量"
// @Param limit query int false "限制数量"
// @Success 200 {object} response.Response
// @Router /api/message/history [get]
func (h *MessageHandler) GetChatHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid friend id")
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, err := h.messageService.GetChatHistory(userID, friendID, offset, limit)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, messages)
}

// GetGroupChatHistory 获取群聊历史
// @Summary 获取群聊历史
// @Tags message
// @Produce json
// @Param group_id query int true "群组ID"
// @Param offset query int false "偏移量"
// @Param limit query int false "限制数量"
// @Success 200 {object} response.Response
// @Router /api/message/group/history [get]
func (h *MessageHandler) GetGroupChatHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	groupIDStr := c.Query("group_id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid group id")
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, err := h.messageService.GetGroupChatHistory(groupID, offset, limit)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, messages)
}

// GetUnreadMessages 获取未读消息
// @Summary 获取未读消息
// @Tags message
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/message/unread [get]
func (h *MessageHandler) GetUnreadMessages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	messages, err := h.messageService.GetUnreadMessages(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, messages)
}

// MarkAsRead 标记消息为已读
// @Summary 标记消息为已读
// @Tags message
// @Accept json
// @Produce json
// @Param request body []int64 true "消息ID列表"
// @Success 200 {object} response.Response
// @Router /api/message/read [post]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var messageIDs []int64
	if err := c.ShouldBindJSON(&messageIDs); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	if err := h.messageService.MarkAsRead(messageIDs); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUnreadCount 获取未读消息数量
// @Summary 获取未读消息数量
// @Tags message
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/message/unread/count [get]
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	count, err := h.messageService.GetUnreadCount(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}

