package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"miniblog/internal/middleware"
	"miniblog/internal/pkg/response"
	"miniblog/internal/service"
	"miniblog/internal/types"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags user
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, user)
}

// GetUserByID 根据 ID 获取用户信息
// @Summary 根据 ID 获取用户信息
// @Tags user
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/user/:id [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.InvalidParam(c, "invalid user id")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param request body types.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req types.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, user)
}

// SearchUsers 搜索用户
// @Summary 搜索用户
// @Tags user
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param offset query int false "偏移量"
// @Param limit query int false "限制数量"
// @Success 200 {object} response.Response
// @Router /api/user/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.InvalidParam(c, "keyword is required")
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, err := h.userService.SearchUsers(keyword, offset, limit)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, users)
}

