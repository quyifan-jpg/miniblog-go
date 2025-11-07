-- 创建数据库
CREATE DATABASE IF NOT EXISTS miniblog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE miniblog;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,      -- bcrypt 加密
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
    remark VARCHAR(100),                 -- 备注名
    status TINYINT DEFAULT 0,            -- 0-待确认 1-已同意 2-已拒绝
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
    group_id BIGINT DEFAULT 0,           -- 0表示单聊
    msg_type TINYINT NOT NULL,           -- 1-文本 2-图片 3-语音 4-视频
    content TEXT NOT NULL,
    status TINYINT DEFAULT 0,            -- 0-未读 1-已读
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_from (from_user_id),
    INDEX idx_to (to_user_id),
    INDEX idx_group (group_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 群组表
CREATE TABLE IF NOT EXISTS `groups` (
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
    role TINYINT DEFAULT 1,              -- 1-普通成员 2-管理员 3-群主
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_group_user (group_id, user_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

