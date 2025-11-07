# MiniBlog IM 快速启动指南

## ✅ 项目已完成配置

- **项目名称**: miniblog
- **数据库名**: miniblog
- **MySQL 密码**: root
- **Redis 密码**: 无（空）

---

## 🚀 方式一：Docker 启动（推荐）

这是最简单的方式，会自动创建数据库、Redis 和应用：

```bash
# 1. 进入项目目录
cd /Users/qyf/Documents/codefield/go-new/miniblog-go

# 2. 启动所有服务
make docker-up

# 等待 15-20 秒让服务完全启动

# 3. 验证服务
curl http://localhost:8080/health
```

### 查看日志
```bash
make docker-logs
```

### 停止服务
```bash
make docker-down
```

---

## 🔧 方式二：本地运行

### 步骤 1: 创建数据库

```bash
# 登录 MySQL
mysql -u root -p

# 输入密码后，执行：
source /Users/qyf/Documents/codefield/go-new/miniblog-go/scripts/init.sql;

# 或者直接执行：
mysql -u root -p < /Users/qyf/Documents/codefield/go-new/miniblog-go/scripts/init.sql
```

### 步骤 2: 确认 Redis 已启动

```bash
# 检查 Redis 是否运行
redis-cli ping

# 如果返回 PONG 说明 Redis 正常
```

### 步骤 3: 修改配置（如需要）

编辑 `config/config.yaml`，确认 MySQL 密码正确：

```yaml
mysql:
  host: "127.0.0.1"
  port: 3306
  database: "miniblog"
  username: "root"
  password: "root"  # 如果你的 MySQL 密码不是 root，请修改这里
```

### 步骤 4: 运行项目

```bash
cd /Users/qyf/Documents/codefield/go-new/miniblog-go

# 方式 A：直接运行
make run

# 方式 B：编译后运行
make build
./bin/miniblog
```

---

## ✅ 验证服务是否成功启动

### 1. 检查服务健康状态

```bash
# 检查 HTTP 服务
curl http://localhost:8080/health

# 期望返回：
# {"service":"http","status":"ok"}

# 检查 WebSocket 服务
curl http://localhost:9090/health

# 期望返回：
# {"online_users":0,"status":"ok"}
```

### 2. 测试注册接口

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "nickname": "测试用户"
  }'
```

成功返回示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "created_at": "2025-11-07T10:00:00Z"
  }
}
```

### 3. 测试登录接口

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

成功返回示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户"
    }
  }
}
```

---

## 📱 测试完整流程

### 使用测试脚本（需要安装 jq）

```bash
# 安装 jq（JSON 处理工具）
# macOS:
brew install jq

# 运行测试脚本
cd /Users/qyf/Documents/codefield/go-new/miniblog-go
./scripts/test-api.sh
```

### 手动测试

详见 `QUICKSTART.md` 文件中的完整测试流程。

---

## 🎯 服务地址

- **HTTP API**: http://localhost:8080
- **WebSocket**: ws://localhost:9090/ws
- **健康检查**: 
  - HTTP: http://localhost:8080/health
  - WebSocket: http://localhost:9090/health

---

## 📖 API 文档

查看 `README.md` 获取完整的 API 文档和使用说明。

---

## ❓ 常见问题

### Q1: Docker 启动失败？

```bash
# 检查端口是否被占用
lsof -i :8080
lsof -i :9090
lsof -i :3306
lsof -i :6379

# 停止占用端口的进程或修改配置文件中的端口
```

### Q2: 连接数据库失败？

**检查清单**：
- [ ] MySQL 是否启动？ `mysql.server status`
- [ ] 数据库 `miniblog` 是否已创建？
- [ ] 配置文件中的密码是否正确？
- [ ] 防火墙是否阻止连接？

### Q3: 连接 Redis 失败？

```bash
# 启动 Redis
redis-server

# 检查 Redis 是否运行
redis-cli ping
```

### Q4: 端口被占用？

修改 `config/config.yaml` 中的端口：

```yaml
server:
  http:
    port: 8081  # 修改为其他可用端口
  websocket:
    port: 9091  # 修改为其他可用端口
```

---

## 📝 查看日志

### Docker 方式
```bash
make docker-logs
```

### 本地运行方式
```bash
tail -f logs/miniblog.log
```

---

## 🛑 停止服务

### Docker 方式
```bash
make docker-down
```

### 本地运行方式
按 `Ctrl + C` 停止进程

---

## 🎉 下一步

1. 查看 `README.md` 了解完整功能
2. 查看 `QUICKSTART.md` 学习详细用法
3. 查看 `PROJECT_OVERVIEW.md` 了解架构设计

---

**祝您使用愉快！** 🚀

