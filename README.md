# Simple IM - 即时通讯系统

一个基于 Go 语言开发的轻量级即时通讯系统，支持单聊、群聊、好友管理等功能。

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端层                              │
│  Web 浏览器 / 移动 App / 桌面应用                             │
└────────────┬───────────────────────────┬────────────────────┘
             │                           │
             │ HTTP API                  │ WebSocket
             ▼                           ▼
┌─────────────────────┐      ┌─────────────────────┐
│   HTTP 服务进程      │      │   Comet 服务进程     │
│   (端口 8080)       │      │   (端口 9090)       │
│                     │      │                     │
│  ┌────────────┐    │      │  ┌────────────┐    │
│  │  Handler   │    │      │  │  Handler   │    │
│  └─────┬──────┘    │      │  └─────┬──────┘    │
│        │           │      │        │           │
│  ┌─────▼──────┐    │      │  ┌─────▼──────┐    │
│  │  Service   │    │      │  │  Subscribe │    │
│  └─────┬──────┘    │      │  └─────┬──────┘    │
│        │           │      │        │           │
│  ┌─────▼──────┐    │      │  ┌─────▼──────┐    │
│  │Repository  │    │      │  │ SessionMgr │    │
│  └────────────┘    │      │  └────────────┘    │
└─────────┬───────────┘      └─────────┬───────────┘
          │                            │
          │                            │
          └────────────┬───────────────┘
                       │
          ┌────────────▼────────────┐
          │                         │
          │    中间件层 (共享)       │
          │                         │
          │  ┌──────┐  ┌────────┐  │
          │  │MySQL │  │ Redis  │  │
          │  │      │  │        │  │
          │  │ 持久化│  │ 缓存/  │  │
          │  │ 存储 │  │ PubSub │  │
          │  └──────┘  └────────┘  │
          └─────────────────────────┘
```

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0
- **WebSocket**: Gorilla WebSocket
- **JWT**: golang-jwt/jwt
- **ORM**: GORM
- **日志**: Zap

## 功能特性

- ✅ 用户注册/登录
- ✅ JWT 认证
- ✅ 好友管理（添加、删除、好友列表）
- ✅ 单聊消息
- ✅ 群聊消息
- ✅ 消息已读/未读
- ✅ WebSocket 实时通信
- ✅ 在线状态管理
- ✅ 心跳检测

## 项目结构

```
simple-im/
├── cmd/                          # 命令行入口
│   └── server/
│       ├── main.go               # 主入口
│       ├── http.go               # HTTP 服务启动
│       └── websocket.go          # WebSocket 服务启动
│
├── config/                       # 配置管理
│   ├── config.go                 # 配置结构定义
│   └── config.yaml               # 配置文件
│
├── internal/                     # 内部代码
│   ├── handler/                  # HTTP 处理层
│   ├── service/                  # 业务逻辑层
│   ├── repository/               # 数据访问层
│   ├── model/                    # 数据模型
│   ├── websocket/                # WebSocket 服务
│   ├── middleware/               # 中间件
│   ├── pkg/                      # 内部公共包
│   └── types/                    # 类型定义
│
├── scripts/                      # 脚本文件
│   └── init.sql                  # 数据库初始化
│
├── docker/                       # Docker 配置
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── go.mod                        # Go 模块定义
├── Makefile                      # 构建脚本
└── README.md                     # 项目文档
```

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

### 本地运行

1. **克隆项目**

```bash
git clone <repository-url>
cd simple-im
```

2. **安装依赖**

```bash
make install
```

3. **配置数据库**

```bash
# 创建数据库
mysql -u root -p < scripts/init.sql
```

4. **修改配置文件**

编辑 `config/config.yaml`，修改数据库和 Redis 连接信息。

5. **运行项目**

```bash
make run
```

### Docker 运行

```bash
# 启动所有服务
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

## API 文档

### 认证相关

#### 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "nickname": "Test User",
  "email": "test@example.com"
}
```

#### 用户登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 用户相关

#### 获取当前用户信息

```http
GET /api/user/profile
Authorization: Bearer <token>
```

#### 更新用户信息

```http
PUT /api/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "New Nickname",
  "avatar": "https://example.com/avatar.jpg"
}
```

#### 搜索用户

```http
GET /api/user/search?keyword=test&offset=0&limit=20
Authorization: Bearer <token>
```

### 好友相关

#### 添加好友

```http
POST /api/friend/add
Authorization: Bearer <token>
Content-Type: application/json

{
  "friend_id": 2,
  "remark": "My Friend"
}
```

#### 处理好友请求

```http
POST /api/friend/handle
Authorization: Bearer <token>
Content-Type: application/json

{
  "friend_id": 2,
  "accept": true
}
```

#### 获取好友列表

```http
GET /api/friend/list
Authorization: Bearer <token>
```

### 消息相关

#### 发送消息

```http
POST /api/message/send
Authorization: Bearer <token>
Content-Type: application/json

{
  "to_user_id": 2,
  "msg_type": 1,
  "content": "Hello!"
}
```

#### 获取聊天历史

```http
GET /api/message/history?friend_id=2&offset=0&limit=50
Authorization: Bearer <token>
```

#### 获取未读消息

```http
GET /api/message/unread
Authorization: Bearer <token>
```

### 群组相关

#### 创建群组

```http
POST /api/group/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My Group",
  "description": "A test group",
  "member_ids": [2, 3, 4]
}
```

#### 获取群组列表

```http
GET /api/group/list
Authorization: Bearer <token>
```

## WebSocket 连接

### 连接

```javascript
const ws = new WebSocket('ws://localhost:9090/ws?token=<your-jwt-token>');

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

### 消息格式

```json
{
  "event": "message",
  "data": {
    "id": 1,
    "from_user_id": 1,
    "to_user_id": 2,
    "msg_type": 1,
    "content": "Hello!",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": 1704067200
}
```

### 事件类型

- `ping` - 心跳检测
- `pong` - 心跳响应
- `message` - 私聊消息
- `group_message` - 群聊消息
- `friend_online` - 好友上线
- `friend_offline` - 好友离线

## 配置说明

```yaml
server:
  http:
    host: "0.0.0.0"
    port: 8080
  websocket:
    host: "0.0.0.0"
    port: 9090

mysql:
  host: "127.0.0.1"
  port: 3306
  database: "simple_im"
  username: "root"
  password: "123456"

redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expire_hours: 168  # 7 天

log:
  level: "info"
  file: "logs/simple-im.log"
```

## 开发指南

### 代码规范

```bash
# 格式化代码
make fmt

# 代码检查
make lint
```

### 测试

```bash
# 运行测试
make test
```

### 编译

```bash
# 编译项目
make build

# 清理构建文件
make clean
```

## 性能优化建议

1. **数据库优化**
   - 为高频查询字段添加索引
   - 使用连接池管理数据库连接
   - 考虑读写分离

2. **缓存优化**
   - 缓存热点数据（用户信息、好友列表）
   - 使用 Redis 实现分布式会话管理

3. **WebSocket 优化**
   - 实现负载均衡
   - 使用 Redis PubSub 实现跨服务器消息推送

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

