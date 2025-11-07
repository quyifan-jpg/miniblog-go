# Simple IM 项目概览

## 项目简介

Simple IM 是一个基于 Go 语言开发的轻量级即时通讯系统，采用前后端分离的架构设计，支持单聊、群聊、好友管理等核心功能。

## 核心特性

### 1. 架构设计

- **分层架构**: Handler -> Service -> Repository 三层架构，职责清晰
- **双服务模式**: HTTP API 服务 (8080) + WebSocket 实时通信服务 (9090)
- **中间件支持**: 认证、跨域、日志等中间件
- **依赖注入**: 通过构造函数注入依赖

### 2. 技术特点

- **高性能**: 基于 Gin 框架，高性能 HTTP 路由
- **实时通信**: WebSocket 支持，实现消息实时推送
- **安全认证**: JWT Token 认证机制
- **数据持久化**: GORM ORM，自动迁移数据表
- **缓存支持**: Redis 缓存 + PubSub 消息分发
- **日志系统**: 基于 Zap 的结构化日志

### 3. 功能模块

#### 用户模块
- 用户注册（密码 bcrypt 加密）
- 用户登录（返回 JWT Token）
- 用户信息管理
- 用户搜索

#### 好友模块
- 添加好友（发送好友请求）
- 处理好友请求（同意/拒绝）
- 好友列表查看
- 好友删除
- 好友备注

#### 消息模块
- 单聊消息发送
- 群聊消息发送
- 聊天历史查询
- 未读消息管理
- 消息已读标记
- 支持多种消息类型（文本、图片、语音、视频）

#### 群组模块
- 创建群组
- 群组信息管理
- 添加/移除群成员
- 群成员角色管理（群主、管理员、普通成员）
- 退出群组

#### WebSocket 模块
- 连接管理
- 会话管理
- 心跳检测
- 在线状态管理
- 实时消息推送

## 项目结构说明

```
simple-im/
├── cmd/server/              # 应用入口
│   ├── main.go             # 主程序入口
│   ├── http.go             # HTTP 服务器启动
│   └── websocket.go        # WebSocket 服务器启动
│
├── config/                  # 配置管理
│   ├── config.go           # 配置结构和加载逻辑
│   └── config.yaml         # 配置文件
│
├── internal/               # 内部代码（不对外暴露）
│   ├── handler/           # HTTP 请求处理层
│   │   ├── auth.go        # 认证处理器
│   │   ├── user.go        # 用户处理器
│   │   ├── friend.go      # 好友处理器
│   │   ├── message.go     # 消息处理器
│   │   └── group.go       # 群组处理器
│   │
│   ├── service/           # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── friend_service.go
│   │   ├── message_service.go
│   │   └── group_service.go
│   │
│   ├── repository/        # 数据访问层
│   │   ├── user_repo.go
│   │   ├── friend_repo.go
│   │   ├── message_repo.go
│   │   └── group_repo.go
│   │
│   ├── model/             # 数据模型
│   │   ├── user.go
│   │   ├── friend.go
│   │   ├── message.go
│   │   └── group.go
│   │
│   ├── websocket/         # WebSocket 服务
│   │   ├── server.go      # WebSocket 服务器
│   │   ├── manager.go     # 连接管理器
│   │   ├── session.go     # 会话管理
│   │   ├── handler.go     # 消息处理
│   │   └── heartbeat.go   # 心跳机制
│   │
│   ├── middleware/        # 中间件
│   │   ├── auth.go        # JWT 认证中间件
│   │   ├── cors.go        # 跨域中间件
│   │   └── logger.go      # 日志中间件
│   │
│   ├── pkg/              # 内部公共包
│   │   ├── jwt/          # JWT 工具
│   │   ├── redis/        # Redis 封装
│   │   ├── mysql/        # MySQL 封装
│   │   ├── logger/       # 日志封装
│   │   └── response/     # 统一响应
│   │
│   └── types/            # 类型定义
│       ├── request.go    # 请求类型
│       ├── response.go   # 响应类型
│       └── event.go      # WebSocket 事件类型
│
├── scripts/              # 脚本文件
│   ├── init.sql         # 数据库初始化脚本
│   └── test-api.sh      # API 测试脚本
│
├── docker/              # Docker 配置
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── go.mod              # Go 模块定义
├── go.sum              # 依赖锁定
├── Makefile            # 构建脚本
├── README.md           # 项目文档
├── QUICKSTART.md       # 快速开始指南
└── PROJECT_OVERVIEW.md # 项目概览（本文件）
```

## 数据库设计

### 用户表 (users)
- id: 主键
- username: 用户名（唯一）
- password: 密码（bcrypt 加密）
- nickname: 昵称
- avatar: 头像
- email: 邮箱
- created_at: 创建时间
- updated_at: 更新时间

### 好友关系表 (friends)
- id: 主键
- user_id: 用户 ID
- friend_id: 好友 ID
- remark: 备注名
- status: 状态（0-待确认, 1-已同意, 2-已拒绝）
- created_at: 创建时间
- updated_at: 更新时间

### 消息表 (messages)
- id: 主键
- from_user_id: 发送者 ID
- to_user_id: 接收者 ID
- group_id: 群组 ID（0 表示单聊）
- msg_type: 消息类型（1-文本, 2-图片, 3-语音, 4-视频）
- content: 消息内容
- status: 状态（0-未读, 1-已读）
- created_at: 创建时间

### 群组表 (groups)
- id: 主键
- name: 群组名称
- avatar: 群组头像
- owner_id: 群主 ID
- description: 群组描述
- created_at: 创建时间
- updated_at: 更新时间

### 群成员表 (group_members)
- id: 主键
- group_id: 群组 ID
- user_id: 用户 ID
- role: 角色（1-普通成员, 2-管理员, 3-群主）
- joined_at: 加入时间

## API 设计

### 认证相关
- POST /api/auth/register - 用户注册
- POST /api/auth/login - 用户登录

### 用户相关
- GET /api/user/profile - 获取当前用户信息
- PUT /api/user/profile - 更新用户信息
- GET /api/user/:id - 根据 ID 获取用户信息
- GET /api/user/search - 搜索用户

### 好友相关
- POST /api/friend/add - 添加好友
- POST /api/friend/handle - 处理好友请求
- GET /api/friend/list - 获取好友列表
- GET /api/friend/requests - 获取好友请求列表
- DELETE /api/friend/:id - 删除好友

### 消息相关
- POST /api/message/send - 发送消息
- GET /api/message/history - 获取聊天历史
- GET /api/message/group/history - 获取群聊历史
- GET /api/message/unread - 获取未读消息
- POST /api/message/read - 标记消息为已读
- GET /api/message/unread/count - 获取未读消息数量

### 群组相关
- POST /api/group/create - 创建群组
- GET /api/group/:id - 获取群组信息
- PUT /api/group/:id - 更新群组信息
- DELETE /api/group/:id - 删除群组
- GET /api/group/list - 获取用户加入的群组列表
- POST /api/group/:id/member/add - 添加群成员
- DELETE /api/group/:id/member/:member_id - 移除群成员
- GET /api/group/:id/members - 获取群成员列表
- POST /api/group/:id/leave - 退出群组

### WebSocket
- GET /ws?token=xxx - 建立 WebSocket 连接

## WebSocket 消息格式

### 基本格式
```json
{
  "event": "message",
  "data": {},
  "msg_id": "optional-message-id",
  "timestamp": 1234567890
}
```

### 事件类型
- `ping`: 心跳检测
- `pong`: 心跳响应
- `message`: 私聊消息
- `group_message`: 群聊消息
- `message_read`: 消息已读
- `friend_online`: 好友上线
- `friend_offline`: 好友离线

## 技术要点

### 1. 认证机制
- 使用 JWT 进行用户认证
- Token 有效期 7 天
- 请求 Header 携带: `Authorization: Bearer <token>`

### 2. 密码安全
- 使用 bcrypt 加密，加密成本系数为默认值
- 不返回密码字段（模型中使用 `json:"-"`）

### 3. WebSocket 连接
- 支持 Token 查询参数或 Header 认证
- 实现心跳检测机制（30秒 ping，60秒超时）
- 自动检测和清理死连接

### 4. 消息推送
- 使用 Redis PubSub 实现跨服务器消息分发
- 异步消息处理，不阻塞主流程
- 消息持久化到 MySQL

### 5. 并发控制
- 使用 sync.RWMutex 保护共享资源
- 会话管理采用线程安全的 Map
- 消息发送采用 channel 缓冲

## 性能优化

### 1. 数据库优化
- 为高频查询字段添加索引
- 使用连接池管理数据库连接
- 批量操作使用事务

### 2. 缓存策略
- 可以缓存用户信息、好友列表等热点数据
- 使用 Redis 作为缓存层

### 3. WebSocket 优化
- 使用 channel 进行消息缓冲
- 独立的读写 goroutine
- 定期清理死连接

## 扩展方向

### 1. 功能扩展
- 文件上传（对接 MinIO）
- 语音/视频通话
- 消息撤回
- 消息转发
- @提及功能
- 表情包支持
- 阅后即焚

### 2. 架构扩展
- 微服务拆分
- 负载均衡
- 读写分离
- 分库分表
- 消息队列（Kafka）

### 3. 安全增强
- HTTPS 支持
- WSS (WebSocket over TLS)
- 限流防刷
- XSS 防护
- SQL 注入防护

### 4. 运维监控
- Prometheus 监控
- 链路追踪
- 日志聚合
- 告警系统

## 开发规范

### 1. 代码规范
- 遵循 Go 官方代码规范
- 使用 golangci-lint 进行代码检查
- 使用 gofmt 格式化代码

### 2. 注释规范
- 公开函数必须有注释
- 复杂逻辑添加必要注释
- 结构体字段添加说明

### 3. 错误处理
- 统一使用 error 返回
- 记录详细的错误日志
- 返回友好的错误信息

### 4. 测试规范
- 编写单元测试
- 测试覆盖率 > 70%
- 集成测试验证主流程

## 部署方案

### 1. 单机部署
```bash
make build
./bin/simple-im
```

### 2. Docker 部署
```bash
make docker-up
```

### 3. Kubernetes 部署
- 创建 Deployment
- 创建 Service
- 配置 Ingress
- 配置 ConfigMap

## 常见问题

### Q1: 如何添加新的 API？
1. 在 `internal/types` 定义请求/响应类型
2. 在 `internal/repository` 添加数据访问方法
3. 在 `internal/service` 添加业务逻辑
4. 在 `internal/handler` 添加 HTTP 处理器
5. 在 `cmd/server/http.go` 注册路由

### Q2: 如何处理 WebSocket 消息？
在 `internal/websocket/session.go` 的 `handleMessage` 方法中添加处理逻辑。

### Q3: 如何添加新的中间件？
1. 在 `internal/middleware` 创建中间件文件
2. 实现 `gin.HandlerFunc` 函数
3. 在路由中使用 `router.Use(middleware)`

### Q4: 如何修改数据库结构？
1. 修改 `internal/model` 中的结构体
2. GORM 会在启动时自动迁移表结构
3. 或者编写 migration 脚本手动迁移

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。

## 许可证

MIT License

