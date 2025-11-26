package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

/*
示例 1：基础聊天
演示如何使用 eino 创建一个简单的 LLM 聊天应用
*/

func main() {
	// 从环境变量获取 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置环境变量 OPENAI_API_KEY")
	}

	ctx := context.Background()

	// 创建 OpenAI ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-3.5-turbo",
		APIKey: apiKey,
		// BaseURL: os.Getenv("OPENAI_BASE_URL"), // 如果使用代理，可以设置
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 失败: %v", err)
	}

	// 准备消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个友好的 AI 助手。",
		},
		{
			Role:    schema.User,
			Content: "你好！请简单介绍一下 Go 语言的特点。",
		},
	}

	fmt.Println("🤖 正在生成回复...")
	fmt.Println()

	// 调用模型生成回复
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成回复失败: %v", err)
	}

	// 输出结果
	fmt.Println("✅ 回复内容:")
	fmt.Println(response.Content)
	fmt.Println()

	// 输出使用的 token 数量
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Printf("📊 Token 使用: 输入=%d, 输出=%d, 总计=%d\n",
			response.ResponseMeta.Usage.PromptTokens,
			response.ResponseMeta.Usage.CompletionTokens,
			response.ResponseMeta.Usage.TotalTokens,
		)
	}
}

/*
运行方式:
export OPENAI_API_KEY=your_api_key
go run 01_basic_chat.go

预期输出:
🤖 正在生成回复...

✅ 回复内容:
Go 语言是由 Google 开发的一门编程语言，具有以下特点：
1. 简洁性 - 语法简单，易于学习
2. 并发性 - 内置 goroutine 和 channel
3. 高性能 - 编译型语言，接近 C 的性能
...

📊 Token 使用: 输入=25, 输出=150, 总计=175
*/

