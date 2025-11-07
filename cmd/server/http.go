package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"miniblog/config"
	"miniblog/internal/handler"
	"miniblog/internal/middleware"
	"miniblog/internal/pkg/jwt"
	"miniblog/internal/pkg/logger"
	"miniblog/internal/repository"
	"miniblog/internal/service"
)

// startHTTPServer 启动 HTTP 服务器
func startHTTPServer(cfg *config.Config, jwtManager *jwt.JWTManager) error {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// 初始化 Repository
	userRepo := repository.NewUserRepository()
	friendRepo := repository.NewFriendRepository()
	messageRepo := repository.NewMessageRepository()
	groupRepo := repository.NewGroupRepository()

	// 初始化 Service
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)
	friendService := service.NewFriendService(friendRepo, userRepo)
	messageService := service.NewMessageService(messageRepo, userRepo, friendRepo, groupRepo)
	groupService := service.NewGroupService(groupRepo, userRepo)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	friendHandler := handler.NewFriendHandler(friendService)
	messageHandler := handler.NewMessageHandler(messageService)
	groupHandler := handler.NewGroupHandler(groupService)

	// 公共路由
	api := router.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(jwtManager))
		{
			// 用户相关
			user := authorized.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.GET("/:id", userHandler.GetUserByID)
				user.GET("/search", userHandler.SearchUsers)
			}

			// 好友相关
			friend := authorized.Group("/friend")
			{
				friend.POST("/add", friendHandler.AddFriend)
				friend.POST("/handle", friendHandler.HandleFriendRequest)
				friend.GET("/list", friendHandler.GetFriendList)
				friend.GET("/requests", friendHandler.GetFriendRequests)
				friend.DELETE("/:id", friendHandler.DeleteFriend)
			}

			// 消息相关
			message := authorized.Group("/message")
			{
				message.POST("/send", messageHandler.SendMessage)
				message.GET("/history", messageHandler.GetChatHistory)
				message.GET("/group/history", messageHandler.GetGroupChatHistory)
				message.GET("/unread", messageHandler.GetUnreadMessages)
				message.POST("/read", messageHandler.MarkAsRead)
				message.GET("/unread/count", messageHandler.GetUnreadCount)
			}

			// 群组相关
			group := authorized.Group("/group")
			{
				group.POST("/create", groupHandler.CreateGroup)
				group.GET("/:id", groupHandler.GetGroup)
				group.PUT("/:id", groupHandler.UpdateGroup)
				group.DELETE("/:id", groupHandler.DeleteGroup)
				group.GET("/list", groupHandler.GetUserGroups)
				group.POST("/:id/member/add", groupHandler.AddGroupMembers)
				group.DELETE("/:id/member/:member_id", groupHandler.RemoveGroupMember)
				group.GET("/:id/members", groupHandler.GetGroupMembers)
				group.POST("/:id/leave", groupHandler.LeaveGroup)
			}
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "http",
		})
	})

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	logger.Info("HTTP server starting", zap.String("addr", addr))

	return router.Run(addr)
}

