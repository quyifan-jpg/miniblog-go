package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"miniblog/config"
	"miniblog/internal/model"
)

var DB *gorm.DB

// InitMySQL 初始化 MySQL 连接
func InitMySQL(cfg *config.MySQLConfig) error {
	dsn := cfg.DSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	DB = db

	// 自动迁移表结构
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}

// autoMigrate 自动迁移表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Friend{},
		&model.Message{},
		&model.Group{},
		&model.GroupMember{},
	)
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

