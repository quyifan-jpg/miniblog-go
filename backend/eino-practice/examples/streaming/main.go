package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

/*
示例 3：流式输出
演示如何实现流式输出，实时显示 AI 生成的内容
*/

func main() {
    // 检查 API Key
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("请设置环境变量 OPENAI_API_KEY")
    }

    ctx := context.Background()

    // 1. 创建 ChatModel（流式输出通过调用 Stream() 方法实现，不需要配置字段）
    chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
        Model:  "gpt-3.5-turbo",
        APIKey: apiKey,
    })
    if err != nil {
        log.Fatalf("创建 ChatModel 失败: %v", err)
    }

    // 2. 准备消息
    messages := []*schema.Message{
        {
            Role:    schema.System,
            Content: "你是一个专业的故事作家。所有回答必须是中文。",
        },
        {
            Role:    schema.User,
            Content: "请写一个关于程序员学习新技术的短故事，大约100字。",
        },
    }

    fmt.Println("📖 AI 正在创作故事...")
    fmt.Println()
    fmt.Print("✨ ")

    // 3. 使用流式生成方法
    stream, err := chatModel.Stream(ctx, messages)
    if err != nil {
        log.Fatalf("流式生成失败: %v", err)
    }

    // 4. 核心：处理流式响应
    // Stream() 返回 *schema.StreamReader[*schema.Message]
    // 每个 chunk 是一个 *schema.Message
    fullContent := ""
    for {
        chunk, err := stream.Recv()
        
        // 流结束标志
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("接收流式响应失败: %v", err)
        }

        // 实时打印内容并累加内容
        if chunk != nil && chunk.Content != "" {
            fmt.Print(chunk.Content)
            fullContent += chunk.Content
        }
    }

    fmt.Println()
    fmt.Println()
    fmt.Printf("✅ 创作完成。总字符数: %d\n", len(fullContent))
}