package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess      = 0
	CodeError        = 1
	CodeInvalidParam = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// InvalidParam 参数错误
func InvalidParam(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeInvalidParam,
		Message: message,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeForbidden,
		Message: message,
	})
}

// NotFound 未找到
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// ServerError 服务器错误
func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeServerError,
		Message: message,
	})
}

