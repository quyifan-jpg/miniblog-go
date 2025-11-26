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

// 模拟外部工具的执行结果
type FinancialData struct {
	StockPrice string
	ReportSummary string
}

func main() {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("请设置环境变量 OPENAI_API_KEY")
    }

    ctx := context.Background()

    // 1. 创建 ChatModel
    chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
        Model:  "gpt-3.5-turbo",
        APIKey: apiKey,
    })
    if err != nil {
        log.Fatalf("创建 ChatModel 失败: %v", err)
    }

    // --- 复杂 Chain 步骤定义 ---
    // 注意：Chain 要求所有步骤的输入/输出类型必须匹配，所以统一使用 string -> string

    // 步骤1: 提取关键实体（string -> string）
    entityExtractor := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
        fmt.Printf("📝 步骤1 - 实体提取：分析用户意图...\n")
        // 实际应用中会调用 LLM 进行实体识别。这里简化为固定实体。
        if input == "Google" {
            return "GOOGL", nil
        }
        return "Unknown_Stock", nil
    })

    // 步骤2: 数据聚合（string -> string，内部处理 FinancialData）
    // 将数据聚合的结果格式化为字符串传递给下一步
    dataGatherer := compose.InvokableLambda(func(ctx context.Context, stockSymbol string) (string, error) {
        fmt.Printf("🔧 步骤2 - 数据聚合：并行查询 (%s)\n", stockSymbol)

        // Tool A: 模拟获取实时价格
        price := "175.50 USD" 
        // Tool B: 模拟获取最新财报摘要
        summary := "近期增长强劲，AI投入高，但短期内面临市场竞争加剧的挑战。"

        // 将 FinancialData 格式化为字符串
        dataStr := fmt.Sprintf("价格:%s|摘要:%s", price, summary)
        return dataStr, nil
    })
    
    // 步骤3: LLM 总结和建议（string -> string）
    // 接收格式化的数据字符串，调用 LLM 生成建议
    summarizer := compose.InvokableLambda(func(ctx context.Context, dataStr string) (string, error) {
        fmt.Printf("🧠 步骤3 - LLM 推理：结合数据生成建议...\n")

        // 解析数据字符串（实际应用中可以使用 JSON）
        // 简化处理：假设格式为 "价格:xxx|摘要:xxx"
        var price, summary string
        fmt.Sscanf(dataStr, "价格:%s|摘要:%s", &price, &summary)

        // 构造给 LLM 的最终 Prompt
        analysisPrompt := fmt.Sprintf(
            "请作为专业金融分析师，根据以下数据给出简短的投资建议：\n价格：%s\n摘要：%s\n建议：",
            price, summary,
        )

        // 调用 LLM
        messages := []*schema.Message{
            {Role: schema.System, Content: "你是一位专业的金融分析师，请根据提供的资料给出买入、卖出或持有建议。"},
            {Role: schema.User, Content: analysisPrompt},
        }

        response, err := chatModel.Generate(ctx, messages)
        if err != nil {
            return "", err
        }
        
        return response.Content, nil
    })

    // 5. 创建 Chain 并添加步骤
    // 所有步骤都是 string -> string，类型匹配
    chain := compose.NewChain[string, string]() 
    chain.AppendLambda(entityExtractor) 
    chain.AppendLambda(dataGatherer) 
    chain.AppendLambda(summarizer)


    // 编译和执行
    runnable, err := chain.Compile(ctx)
    if err != nil {
        log.Fatalf("Chain 编译失败: %v", err)
    }

    fmt.Println("▶️  开始执行复杂金融分析 Chain...")
    fmt.Println()

    result, err := runnable.Invoke(ctx, "Google 股票怎么样？")
    if err != nil {
        log.Fatalf("Chain 执行失败: %v", err)
    }

    fmt.Println()
    fmt.Println("✅ 最终分析结果:")
    fmt.Println(result)
}
// ⚠️ 注意：NewAdapter 是为了演示目的而假设 Eino 框架提供的工具，用于连接不同输入/输出类型的 Lambda。