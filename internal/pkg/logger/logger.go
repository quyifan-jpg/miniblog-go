package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger 初始化日志
func InitLogger(level, logFile string, maxSize, maxBackups, maxAge int, compress bool) error {
	// 确保日志目录存在
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 解析日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建 Core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		),
		zapLevel,
	)

	// 创建 Logger
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Info 记录信息日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Debug 记录调试日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Warn 记录警告日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error 记录错误日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal 记录致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// Sync 同步日志
func Sync() error {
	return Log.Sync()
}

