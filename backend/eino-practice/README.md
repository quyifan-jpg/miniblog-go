# Eino Go 框架学习实践

## 📚 关于 Eino

Eino 是字节跳动开源的 LLM Application 开发框架，专注于复杂 AI 应用的开发。

### 特性
- 🎯 **组件化设计** - ChatModel、Tool、Memory、Retriever 等可复用组件
- 🔗 **编排能力** - Chain、Graph 等编排方式
- 🌊 **流式处理** - 原生支持流式输出
- 🔌 **扩展性强** - 易于集成各种 LLM 和工具

### 官方资源
- GitHub: https://github.com/cloudwego/eino
- 文档: https://www.cloudwego.io/zh/docs/eino/

---

## 🎯 学习路线

### 第一阶段：基础使用
1. ✅ 基础聊天 - `examples/basic-chat/`
2. ✅ 带工具的 Agent - `examples/agent-with-tools/`
3. ✅ 流式输出 - `examples/streaming/`

### 第二阶段：进阶功能
4. ✅ Chain 编排 - `examples/chain/`
5. ⬜ Memory 记忆 - `examples/memory/`
6. ⬜ RAG 检索 - `examples/rag/`

### 第三阶段：复杂应用
7. ⬜ Graph 工作流 - `examples/graph/`
8. ⬜ 自定义组件 - `examples/custom-component/`

---

## 🚀 快速开始

### 1. 安装依赖

```bash
cd backend/eino-practice
go mod tidy
```

### 2. 配置 API Key

```bash
# 创建 .env 文件
cat > .env << EOF
OPENAI_API_KEY=your_api_key_here
OPENAI_BASE_URL=https://api.openai.com/v1
EOF
```

### 3. 运行示例

```bash
# 基础聊天
go run examples/basic-chat/main.go

# Agent with Tools
go run examples/agent-with-tools/main.go

# 流式输出
go run examples/streaming/main.go

# Chain 编排
go run examples/chain/main.go
```

---

## 📖 核心概念

### 1. ChatModel（聊天模型）
```go
model, _ := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    Model: "gpt-3.5-turbo",
})
```

### 2. Tool（工具）
```go
tool := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
    return "result", nil
})
```

### 3. Chain（链式调用）
```go
chain := compose.NewChain[string, string]()
chain.AppendChatModel(model)
```

### 4. Graph（图工作流）
```go
graph := compose.NewGraph[State]()
graph.AddNode("node1", node1Func)
graph.AddNode("node2", node2Func)
graph.AddEdge("node1", "node2")
```

---

## 🔧 目录结构

```
backend/eino-practice/
├── go.mod                      # Go 模块文件
├── go.sum                      # 依赖锁定文件
├── README.md                   # 本文件
├── QUICKSTART.md               # 快速开始指南
├── config.example              # 配置示例
└── examples/                   # 示例代码目录
    ├── basic-chat/            # 基础聊天示例
    │   └── main.go
    ├── agent-with-tools/      # Agent 工具示例
    │   └── main.go
    ├── streaming/             # 流式输出示例
    │   └── main.go
    └── chain/                 # Chain 编排示例
        └── main.go
```

---

## 🎓 学习建议

1. **按顺序学习** - 从简单到复杂
2. **动手实践** - 修改示例代码，观察效果
3. **查看源码** - 理解框架设计思路
4. **参考文档** - 官方文档有详细说明

---

## 💡 常见问题

### Q: 支持哪些 LLM？
A: OpenAI、Anthropic、Azure OpenAI、本地模型等

### Q: 如何使用国内 API？
A: 修改 `OPENAI_BASE_URL` 为国内代理地址

### Q: 如何调试？
A: 使用 `compose.WithDebugMode()` 开启调试模式

---

## 📚 参考资料

- [Eino GitHub](https://github.com/cloudwego/eino)
- [官方文档](https://www.cloudwego.io/zh/docs/eino/)
- [示例代码](https://github.com/cloudwego/eino/tree/main/examples)

