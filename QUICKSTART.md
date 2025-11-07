# 快速开始指南

## 环境准备

### 1. 安装依赖

确保已安装以下软件：

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

### 2. 初始化数据库

```bash
# 登录 MySQL
mysql -u root -p

# 执行初始化脚本
source scripts/init.sql;
```

### 3. 修改配置

编辑 `config/config.yaml` 文件，修改数据库和 Redis 连接信息：

```yaml
mysql:
  host: "127.0.0.1"
  port: 3306
  database: "simple_im"
  username: "root"
  password: "your_password"  # 修改为你的密码

redis:
  host: "127.0.0.1"
  port: 6379
  password: ""  # 如果有密码，在这里填写
```

## 运行方式

### 方式一：直接运行

```bash
# 下载依赖
go mod tidy

# 运行项目
go run cmd/server/main.go cmd/server/http.go cmd/server/websocket.go
```

### 方式二：使用 Makefile

```bash
# 安装依赖
make install

# 运行项目
make run
```

### 方式三：编译后运行

```bash
# 编译
make build

# 运行
./bin/simple-im
```

### 方式四：使用 Docker

```bash
# 启动所有服务（包括 MySQL 和 Redis）
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

## 验证服务

### 1. 检查服务状态

```bash
# HTTP 服务健康检查
curl http://localhost:8080/health

# WebSocket 服务健康检查
curl http://localhost:9090/health
```

### 2. 测试注册功能

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "nickname": "测试用户",
    "email": "test@example.com"
  }'
```

### 3. 测试登录功能

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

返回结果会包含 token，保存这个 token 用于后续请求。

### 4. 测试获取用户信息

```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## WebSocket 测试

### 使用浏览器控制台测试

打开浏览器控制台（F12），执行以下代码：

```javascript
// 替换为你的 token
const token = 'YOUR_TOKEN_HERE';

// 建立 WebSocket 连接
const ws = new WebSocket(`ws://localhost:9090/ws?token=${token}`);

ws.onopen = () => {
  console.log('✅ WebSocket 连接成功');
  
  // 发送 ping
  ws.send(JSON.stringify({
    event: 'ping',
    timestamp: Date.now()
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('📨 收到消息:', message);
};

ws.onerror = (error) => {
  console.error('❌ WebSocket 错误:', error);
};

ws.onclose = () => {
  console.log('🔌 WebSocket 连接关闭');
};
```

## 完整测试流程

### 1. 创建两个测试账号

```bash
# 用户 A
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123",
    "nickname": "Alice"
  }'

# 用户 B
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "password": "password123",
    "nickname": "Bob"
  }'
```

### 2. 登录获取 token

```bash
# Alice 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123"
  }'

# Bob 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "password": "password123"
  }'
```

### 3. Alice 添加 Bob 为好友

```bash
# 假设 Bob 的 ID 是 2
curl -X POST http://localhost:8080/api/friend/add \
  -H "Authorization: Bearer ALICE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "friend_id": 2,
    "remark": "我的朋友 Bob"
  }'
```

### 4. Bob 查看好友请求

```bash
curl -X GET http://localhost:8080/api/friend/requests \
  -H "Authorization: Bearer BOB_TOKEN"
```

### 5. Bob 接受好友请求

```bash
# 假设 Alice 的 ID 是 1
curl -X POST http://localhost:8080/api/friend/handle \
  -H "Authorization: Bearer BOB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "friend_id": 1,
    "accept": true
  }'
```

### 6. Alice 发送消息给 Bob

```bash
curl -X POST http://localhost:8080/api/message/send \
  -H "Authorization: Bearer ALICE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to_user_id": 2,
    "msg_type": 1,
    "content": "Hi Bob! 你好！"
  }'
```

### 7. Bob 查看消息

```bash
# 查看聊天历史
curl -X GET "http://localhost:8080/api/message/history?friend_id=1&offset=0&limit=20" \
  -H "Authorization: Bearer BOB_TOKEN"

# 查看未读消息
curl -X GET http://localhost:8080/api/message/unread \
  -H "Authorization: Bearer BOB_TOKEN"
```

## 常见问题

### 1. 连接数据库失败

检查：
- MySQL 是否启动
- 配置文件中的数据库信息是否正确
- 数据库是否已创建

### 2. 连接 Redis 失败

检查：
- Redis 是否启动
- 配置文件中的 Redis 信息是否正确

### 3. WebSocket 连接失败

检查：
- WebSocket 服务是否启动（端口 9090）
- Token 是否有效
- 浏览器是否支持 WebSocket

### 4. 端口被占用

修改 `config/config.yaml` 中的端口配置：

```yaml
server:
  http:
    port: 8081  # 修改为其他端口
  websocket:
    port: 9091  # 修改为其他端口
```

## 日志查看

日志文件位于 `logs/simple-im.log`

```bash
# 实时查看日志
tail -f logs/simple-im.log
```

## 停止服务

如果是直接运行，按 `Ctrl+C` 停止。

如果是 Docker 运行：

```bash
make docker-down
```

## 下一步

- 查看 [README.md](README.md) 了解完整的 API 文档
- 查看 [API 文档](README.md#api-文档) 了解所有可用接口
- 查看 [WebSocket 文档](README.md#websocket-连接) 了解实时通信

