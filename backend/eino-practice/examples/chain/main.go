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

/*
示例 4：Chain 编排
演示如何使用 Chain 将多个步骤串联起来
*/

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
		log.Fatalf("创建 ChatModel 失败: %v", err)
	}

	fmt.Println("🔗 Chain 编排示例")
	fmt.Println()

	// 创建一个简单的 Chain
	// 
	// Chain 数据流说明：
	// 1. 初始输入: runnable.Invoke(ctx, "什么是 Eino 框架？") 
	//    ↓
	// 2. preprocessor 的 input = "什么是 Eino 框架？" (原始用户输入)
	//    ↓ 输出: "用户问题: 什么是 Eino 框架？"
	// 3. llmStep 的 input = "用户问题: 什么是 Eino 框架？" (preprocessor 的输出)
	//    ↓ 输出: LLM 生成的回复内容
	// 4. postprocessor 的 input = LLM 的回复内容 (llmStep 的输出)
	//    ↓ 输出: "【AI 回复】\n{LLM回复内容}"
	
	// 步骤1: 预处理输入
	// input: 原始用户输入 (例如: "什么是 Eino 框架？")
	preprocessor := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		fmt.Printf("📝 步骤1 - 预处理: 清理输入...\n")
		fmt.Printf("   [输入] input = %q\n", input)
		cleaned := fmt.Sprintf("用户问题: %s", input)
		fmt.Printf("   [输出] cleaned = %q\n", cleaned)
		return cleaned, nil
	})

	// 步骤2: 调用 LLM
	// input: preprocessor 的输出 (例如: "用户问题: 什么是 Eino 框架？")
	llmStep := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		fmt.Printf("🤖 步骤2 - 调用 LLM: 生成回答...\n")
		fmt.Printf("   [输入] input = %q\n", input)
	
		// 1. 应用 Prompt 模板
		promptTemplate := "请回答以下问题：%s"
		finalPrompt := fmt.Sprintf(promptTemplate, input)
		fmt.Printf("   [处理] finalPrompt = %q\n", finalPrompt)
	
		// 2. 实际调用 LLM
		messages := []*schema.Message{
			{Role: schema.User, Content: finalPrompt},
		}
		
		response, err := chatModel.Generate(ctx, messages) 
		if err != nil {
			return "", err
		}
	
		fmt.Printf("   [输出] response.Content = %q\n", response.Content)
		return response.Content, nil
	})

	// 步骤3: 后处理输出
	// input: llmStep 的输出 (例如: LLM 生成的回复内容)
	postprocessor := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		fmt.Printf("✨ 步骤3 - 后处理: 格式化输出...\n")
		fmt.Printf("   [输入] input = %q\n", input)
		formatted := fmt.Sprintf("【AI 回复】\n%s", input)
		fmt.Printf("   [输出] formatted = %q\n", formatted)
		return formatted, nil
	})

	// 创建 Chain 并添加步骤
	chain := compose.NewChain[string, string]()
	chain.AppendLambda(preprocessor)
	chain.AppendLambda(llmStep)
	chain.AppendLambda(postprocessor)

	// 编译 Chain
	fmt.Println("▶️  开始执行 Chain...")
	fmt.Println()

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("Chain 编译失败: %v", err)
	}

	// 执行 Chain
	result, err := runnable.Invoke(ctx, "什么是 Eino 框架？")
	if err != nil {
		log.Fatalf("Chain 执行失败: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ 最终结果:")
	fmt.Println(result)
}

/*
运行方式:
export OPENAI_API_KEY=your_api_key
go run 04_chain.go

预期输出:
🔗 Chain 编排示例

▶️  开始执行 Chain...

📝 步骤1 - 预处理: 清理输入...
🤖 步骤2 - 调用 LLM: 生成回答...
✨ 步骤3 - 后处理: 格式化输出...

✅ 最终结果:
【AI 回复】
这是 AI 生成的回答
*/

