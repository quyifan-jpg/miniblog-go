package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"miniblog/config"
	"miniblog/internal/pkg/jwt"
	"miniblog/internal/pkg/logger"
	"miniblog/internal/pkg/mysql"
	"miniblog/internal/pkg/redis"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config/config.yaml", "config file path")
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(
		cfg.Log.Level,
		cfg.Log.File,
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
		cfg.Log.Compress,
	); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Simple IM starting...")

	// 初始化 MySQL
	if err := mysql.InitMySQL(&cfg.MySQL); err != nil {
		logger.Fatal("Failed to init MySQL", zap.Error(err))
	}
	defer mysql.Close()
	logger.Info("MySQL initialized")

	// 初始化 Redis
	if err := redis.InitRedis(&cfg.Redis); err != nil {
		logger.Fatal("Failed to init Redis", zap.Error(err))
	}
	defer redis.Close()
	logger.Info("Redis initialized")

	// 创建 JWT Manager
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 启动 WebSocket 服务器（在独立的 goroutine 中）
	go func() {
		if err := startWebSocketServer(cfg, jwtManager); err != nil {
			logger.Error("WebSocket server error", zap.Error(err))
		}
	}()

	// 启动 HTTP 服务器（在独立的 goroutine 中）
	go func() {
		if err := startHTTPServer(cfg, jwtManager); err != nil {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")
	logger.Info("Servers stopped")
}

