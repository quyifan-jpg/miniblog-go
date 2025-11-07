package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	JWT    JWTConfig    `yaml:"jwt"`
	Log    LogConfig    `yaml:"log"`
}

type ServerConfig struct {
	HTTP      HTTPConfig      `yaml:"http"`
	WebSocket WebSocketConfig `yaml:"websocket"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type WebSocketConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config file error: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// DSN 返回 MySQL 连接字符串
func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database)
}

// Addr 返回 Redis 连接地址
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

