# 🚀 Eino Go 快速开始指南

## 📦 安装依赖

```bash
cd backend/eino-practice

# 下载依赖
go mod tidy
```

## 🔑 配置 API Key

### 方式一：使用环境变量（推荐）

```bash
# macOS/Linux
export OPENAI_API_KEY="sk-your-api-key-here"

# Windows PowerShell
$env:OPENAI_API_KEY="sk-your-api-key-here"

# Windows CMD
set OPENAI_API_KEY=sk-your-api-key-here
```

### 方式二：代码中直接设置（不推荐）

```go
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    Model:  "gpt-3.5-turbo",
    APIKey: "sk-your-api-key-here", // 不安全！
})
```

## 🎯 运行示例

### 1️⃣ 基础聊天

```bash
# 设置 API Key
export OPENAI_API_KEY="your_key"

# 运行
go run examples/basic-chat/main.go
```

**预期输出：**
```
🤖 正在生成回复...

✅ 回复内容:
Go 语言是由 Google 开发的编程语言...

📊 Token 使用: 输入=25, 输出=150, 总计=175
```

### 2️⃣ Agent with Tools

```bash
go run examples/agent-with-tools/main.go
```

**预期输出：**
```
🤖 Agent 开始工作...

📍 模型决定调用工具:
   - get_current_time

🔧 调用工具: get_current_time
   参数: map[timezone:Asia/Shanghai]

✅ 最终回复:
现在是 2025-11-07 16:30:45
```

### 3️⃣ 流式输出

```bash
go run examples/streaming/main.go
```

**预期输出：**
```
📖 AI 正在创作故事...

✨ 小李是一名后端工程师，最近听说了 Eino 框架...
（内容会逐字显示）

📊 总字符数: 156
```

### 4️⃣ Chain 编排

```bash
go run examples/chain/main.go
```

**预期输出：**
```
🔗 Chain 编排示例

▶️  开始执行 Chain...

📝 步骤1 - 预处理: 清理输入...
🤖 步骤2 - 调用 LLM: 生成回答...
✨ 步骤3 - 后处理: 格式化输出...

✅ 最终结果:
【AI 回复】
这是 AI 生成的回答
```

## 🎓 学习路径

### 初学者
1. ✅ 先运行 `examples/basic-chat/` 了解基本用法
2. ✅ 再运行 `examples/agent-with-tools/` 学习工具调用
3. ✅ 然后运行 `examples/streaming/` 了解流式输出

### 进阶者
4. ✅ 学习 `examples/chain/` Chain 编排
5. ⬜ 学习 `examples/memory/` 记忆管理
6. ⬜ 学习 `examples/rag/` RAG 检索

### 高级者
7. ⬜ 学习 `examples/graph/` Graph 工作流
8. ⬜ 学习 `examples/custom-component/` 自定义组件

## 💡 常见问题

### Q1: API Key 从哪里获取？

**OpenAI:**
- 访问 https://platform.openai.com/api-keys
- 创建新的 API Key

**国内替代方案:**
- 使用国内代理服务
- 使用开源大模型（如通义千问、文心一言等）

### Q2: 如何使用国内 API？

修改 BaseURL：

```go
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    Model:   "gpt-3.5-turbo",
    APIKey:  apiKey,
    BaseURL: "https://your-proxy.com/v1", // 修改这里
})
```

### Q3: 运行报错 "connection refused"？

可能原因：
1. 网络问题 - 需要代理访问 OpenAI
2. API Key 错误
3. BaseURL 配置错误

解决方案：
```bash
# 使用代理（macOS/Linux）
export HTTPS_PROXY=http://127.0.0.1:7890

# 或修改 BaseURL 使用国内代理
```

### Q4: Token 使用量太大？

优化建议：
1. 使用更便宜的模型（gpt-3.5-turbo）
2. 减少 System Prompt 长度
3. 限制最大 Token 数量

```go
chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    Model:     "gpt-3.5-turbo",
    MaxTokens: 500, // 限制最大输出
})
```

## 🔧 调试技巧

### 1. 开启详细日志

```go
import "github.com/cloudwego/eino/callbacks"

// 创建带回调的模型
handler := callbacks.NewHandlerBuilder().
    OnStart(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
        fmt.Printf("开始: %s\n", info.Name)
        return ctx
    }).
    OnEnd(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
        fmt.Printf("结束: %s\n", info.Name)
        return ctx
    }).
    Build()

response, err := chatModel.Generate(ctx, messages, model.WithCallbacks(handler))
```

### 2. 查看原始请求

```go
// 打印消息内容
for i, msg := range messages {
    fmt.Printf("消息 %d [%s]: %s\n", i, msg.Role, msg.Content)
}
```

### 3. 测试连接

创建简单的测试文件 `test_connection.go`：

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/cloudwego/eino/components/model"
)

func main() {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("请设置 OPENAI_API_KEY")
    }
    
    ctx := context.Background()
    chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
        Model:  "gpt-3.5-turbo",
        APIKey: apiKey,
    })
    if err != nil {
        log.Fatalf("创建失败: %v", err)
    }
    
    messages := []*model.Message{
        {Role: model.User, Content: "Hi!"},
    }
    
    response, err := chatModel.Generate(ctx, messages)
    if err != nil {
        log.Fatalf("调用失败: %v", err)
    }
    
    fmt.Printf("✅ 连接成功! 回复: %s\n", response.Content)
}
```

## 📚 更多资源

- **官方文档**: https://www.cloudwego.io/zh/docs/eino/
- **GitHub**: https://github.com/cloudwego/eino
- **示例代码**: https://github.com/cloudwego/eino/tree/main/examples
- **社区讨论**: https://github.com/cloudwego/eino/discussions

## 🎉 下一步

完成基础示例后，你可以：

1. **集成到 MiniBlog**
   - 添加 AI 聊天功能
   - 实现智能回复
   - 内容审核

2. **探索高级功能**
   - Graph 工作流
   - Memory 管理
   - RAG 检索

3. **构建实际应用**
   - AI 客服
   - 内容生成
   - 智能助手

祝学习愉快！🚀

