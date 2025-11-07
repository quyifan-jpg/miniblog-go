#!/bin/bash

# MiniBlog 数据库创建脚本
# 使用方法: ./scripts/setup-db.sh

set -e  # 遇到错误立即退出

echo "========================================="
echo "MiniBlog 数据库初始化脚本"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_NAME="miniblog"

# 提示输入 MySQL 密码
echo -e "${YELLOW}请输入 MySQL ${DB_USER} 用户的密码:${NC}"
read -s DB_PASSWORD
echo ""

echo -e "${BLUE}连接信息:${NC}"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""

# 测试 MySQL 连接
echo -e "${YELLOW}正在测试 MySQL 连接...${NC}"
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}❌ 无法连接到 MySQL！${NC}"
    echo -e "${YELLOW}请检查:${NC}"
    echo "  1. MySQL 是否已启动"
    echo "  2. 用户名和密码是否正确"
    echo "  3. 主机和端口是否正确"
    exit 1
fi
echo -e "${GREEN}✓ MySQL 连接成功${NC}"
echo ""

# 创建数据库和表
echo -e "${YELLOW}正在创建数据库和表...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" << 'EOSQL'

-- 创建数据库
CREATE DATABASE IF NOT EXISTS miniblog 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE miniblog;

-- 删除已存在的表（谨慎使用）
-- DROP TABLE IF EXISTS group_members;
-- DROP TABLE IF EXISTS messages;
-- DROP TABLE IF EXISTS friends;
-- DROP TABLE IF EXISTS `groups`;
-- DROP TABLE IF EXISTS users;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL COMMENT 'bcrypt 加密',
    nickname VARCHAR(100) DEFAULT NULL,
    avatar VARCHAR(255) DEFAULT NULL,
    email VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
COMMENT='用户表';

-- 好友关系表
CREATE TABLE IF NOT EXISTS friends (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    friend_id BIGINT NOT NULL,
    remark VARCHAR(100) DEFAULT NULL COMMENT '备注名',
    status TINYINT DEFAULT 0 COMMENT '0-待确认 1-已同意 2-已拒绝',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_friend (user_id, friend_id),
    INDEX idx_friend_id (friend_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='好友关系表';

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT NOT NULL,
    group_id BIGINT DEFAULT 0 COMMENT '0表示单聊',
    msg_type TINYINT NOT NULL COMMENT '1-文本 2-图片 3-语音 4-视频',
    content TEXT NOT NULL,
    status TINYINT DEFAULT 0 COMMENT '0-未读 1-已读',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_from (from_user_id),
    INDEX idx_to (to_user_id),
    INDEX idx_group (group_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='消息表';

-- 群组表
CREATE TABLE IF NOT EXISTS `groups` (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(255) DEFAULT NULL,
    owner_id BIGINT NOT NULL,
    description VARCHAR(500) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner (owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='群组表';

-- 群成员表
CREATE TABLE IF NOT EXISTS group_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role TINYINT DEFAULT 1 COMMENT '1-普通成员 2-管理员 3-群主',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_group_user (group_id, user_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='群成员表';

-- 显示创建的表
SHOW TABLES;

EOSQL

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}✅ 数据库创建成功！${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo -e "${BLUE}数据库信息:${NC}"
    echo -e "  数据库名: ${GREEN}miniblog${NC}"
    echo ""
    echo -e "${BLUE}已创建的表:${NC}"
    echo "  ✓ users          - 用户表"
    echo "  ✓ friends        - 好友关系表"
    echo "  ✓ messages       - 消息表"
    echo "  ✓ groups         - 群组表"
    echo "  ✓ group_members  - 群成员表"
    echo ""
    echo -e "${YELLOW}下一步操作:${NC}"
    echo "  1. 确认 config/config.yaml 中的数据库配置:"
    echo "     - database: miniblog"
    echo "     - username: $DB_USER"
    echo "     - password: 你的密码"
    echo ""
    echo "  2. 运行项目:"
    echo -e "     ${GREEN}make run${NC}"
    echo ""
    echo "  3. 验证服务:"
    echo -e "     ${GREEN}curl http://localhost:8080/health${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}=========================================${NC}"
    echo -e "${RED}❌ 数据库创建失败${NC}"
    echo -e "${RED}=========================================${NC}"
    echo ""
    echo -e "${YELLOW}请检查:${NC}"
    echo "  1. MySQL 是否正常运行"
    echo "  2. 用户是否有创建数据库的权限"
    echo "  3. 查看上面的错误信息"
    echo ""
    exit 1
fi

