#!/bin/bash

# MiniBlog 数据库初始化脚本

echo "========================================="
echo "MiniBlog 数据库初始化"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 提示输入 MySQL root 密码
echo -e "${YELLOW}请输入 MySQL root 密码:${NC}"
read -s MYSQL_PASSWORD

echo ""
echo -e "${GREEN}开始创建数据库和表...${NC}"
echo ""

# 执行 SQL 脚本
mysql -u root -p${MYSQL_PASSWORD} << EOF
-- 创建数据库
CREATE DATABASE IF NOT EXISTS miniblog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE miniblog;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar VARCHAR(255),
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 好友关系表
CREATE TABLE IF NOT EXISTS friends (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    friend_id BIGINT NOT NULL,
    remark VARCHAR(100),
    status TINYINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_friend (user_id, friend_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT NOT NULL,
    group_id BIGINT DEFAULT 0,
    msg_type TINYINT NOT NULL,
    content TEXT NOT NULL,
    status TINYINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_from (from_user_id),
    INDEX idx_to (to_user_id),
    INDEX idx_group (group_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 群组表
CREATE TABLE IF NOT EXISTS \`groups\` (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(255),
    owner_id BIGINT NOT NULL,
    description VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner (owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 群成员表
CREATE TABLE IF NOT EXISTS group_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role TINYINT DEFAULT 1,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_group_user (group_id, user_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

EOF

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}✅ 数据库创建成功！${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo -e "数据库名: ${GREEN}miniblog${NC}"
    echo -e "已创建的表:"
    echo -e "  - users (用户表)"
    echo -e "  - friends (好友关系表)"
    echo -e "  - messages (消息表)"
    echo -e "  - groups (群组表)"
    echo -e "  - group_members (群成员表)"
    echo ""
    echo -e "${YELLOW}下一步:${NC}"
    echo -e "  1. 修改 config/config.yaml 中的数据库密码"
    echo -e "  2. 运行: ${GREEN}make run${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}=========================================${NC}"
    echo -e "${RED}❌ 数据库创建失败${NC}"
    echo -e "${RED}=========================================${NC}"
    echo ""
    echo -e "${YELLOW}请检查:${NC}"
    echo -e "  1. MySQL 是否已启动"
    echo -e "  2. root 密码是否正确"
    echo -e "  3. 是否有创建数据库的权限"
    echo ""
fi

