# Backend 练习目录

这个目录用于学习和实践各种 Go 框架和技术。

## 📁 目录结构

```
backend/
└── eino-practice/          # Eino Go 框架学习
    ├── README.md           # Eino 学习指南
    ├── QUICKSTART.md       # 快速开始
    ├── 01_basic_chat.go    # ✅ 基础聊天示例
    ├── 02_agent_with_tools.go  # ✅ Agent 工具示例
    ├── 03_streaming.go     # ✅ 流式输出示例
    └── 04_chain.go         # ✅ Chain 编排示例
```

## 🎯 当前可用练习

### 1. Eino Go 框架（AI Agent 开发）

**位置**: `backend/eino-practice/`

**已准备的示例**:
- ✅ 基础聊天 - `01_basic_chat.go`
- ✅ Agent工具 - `02_agent_with_tools.go`
- ✅ 流式输出 - `03_streaming.go`
- ✅ Chain编排 - `04_chain.go`

**快速开始**:
```bash
cd backend/eino-practice

# 设置 API Key
export OPENAI_API_KEY="your_key"

# 运行示例
go run 01_basic_chat.go
```

**详细文档**: 查看 `eino-practice/README.md` 和 `QUICKSTART.md`

---

## 🚀 如何添加新的练习

### 1. 创建新目录

```bash
cd backend
mkdir your-practice-name
cd your-practice-name
go mod init your-practice-name
```

### 2. 创建 README

```bash
cat > README.md << 'EOF'
# Your Practice Name

## 学习目标
...

## 示例代码
...
EOF
```

### 3. 添加示例代码

创建您的 Go 文件并开始练习！

---

## 📚 推荐学习路径

1. **Eino Go** - AI Agent 开发框架
   - 适合：AI 应用开发
   - 难度：⭐⭐⭐
   - 时间：2-3天

2. **其他框架**（待添加）
   - gRPC
   - GraphQL
   - 微服务框架（Go-Micro, Kratos等）

---

## 💡 学习建议

1. **动手实践** - 运行每个示例并修改代码
2. **查看源码** - 理解框架设计思路
3. **做笔记** - 记录遇到的问题和解决方案
4. **做项目** - 将学到的知识应用到实际项目

---

## 🔗 相关资源

- [Eino GitHub](https://github.com/cloudwego/eino)
- [Eino 文档](https://www.cloudwego.io/zh/docs/eino/)
- [Go 官方文档](https://go.dev/doc/)

---

**Happy Coding!** 🎉

