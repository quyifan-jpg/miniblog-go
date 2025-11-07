package handler

import (
	"github.com/gin-gonic/gin"

	"miniblog/internal/pkg/response"
	"miniblog/internal/service"
	"miniblog/internal/types"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.LoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParam(c, err.Error())
		return
	}

	loginResp, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, loginResp)
}

