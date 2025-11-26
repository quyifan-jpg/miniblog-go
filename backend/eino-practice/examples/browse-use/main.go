package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置环境变量 OPENAI_API_KEY")
	}

	ctx := context.Background()

	// 创建 ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-3.5-turbo",
		APIKey: apiKey,
	})
	if err != nil {
		fmt.Print(err)
	}
	

}