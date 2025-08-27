package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chiatzenw-cur/resume-mixer/internal/tools"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 创建ChatModel
	temp := float32(0.7)
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:       "deepseek-chat",
		APIKey:      os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL:     "https://api.deepseek.com/v1",
		Temperature: &temp,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建工具
	taskTools := []tool.BaseTool{
		&tools.ReadExperienceTool{},
		&tools.ReadSkillsTool{},
		&tools.ReadProjectsTool{},
		&tools.ReadEducationTool{},
		&tools.ReadExtrasTool{},
		&tools.ComposeResumeTool{},
	}

	// 自定义StreamToolCallChecker，用于检测DeepSeek模型的工具调用
	deepSeekStreamToolCallChecker := func(_ context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
		defer sr.Close()

		// 读取整个流以检测工具调用
		// DeepSeek模型可能会在流的任何位置输出工具调用
		for {
			msg, err := sr.Recv()
			if err == io.EOF {
				return false, nil
			}
			if err != nil {
				return false, err
			}

			// 检查是否有工具调用
			if len(msg.ToolCalls) > 0 {
				return true, nil
			}
			if msg != nil && msg.Content != "" {
				fmt.Print(msg.Content)
				//assistantContent.WriteString(chunk.Content)
			}

			// 对于DeepSeek模型，我们继续读取直到流结束或找到工具调用
			// 不像Claude模型那样需要特殊处理，DeepSeek通常在开始就输出工具调用
		}
	}

	// 创建React Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: taskTools,
		},
		StreamToolCallChecker: deepSeekStreamToolCallChecker,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化消息历史
	messages := []*schema.Message{
		schema.SystemMessage(`你是一个专业的简历分析助手。请按照以下步骤进行思考和回应：

1. 当用户提供职位描述(JD)时，立即调用所有工具来读取候选人的简历信息：
   - read_experience: 读取工作经验
   - read_skills: 读取技能信息
   - read_projects: 读取项目经验
   - read_education: 读取教育背景
   - read_extras: 读取其他信息
2. 分析职位描述的技术要求和职责
3. 基于获取的简历信息和职位要求，分析匹配度并提供：
   - 匹配度评分(0-100分)
   - 优势点分析
   - 不足点识别
   - 改进建议
   - 项目想法
   - 推荐书籍
   - 简历和面试突出建议
4. when asked 用compose_resume工具生成最终的简历报告，该工具需要以下参数：
   - job_description: 职位描述（直接使用用户提供的原始JD）
   - relevant_projects: 与职位最相关的项目（需要你分析后筛选）
   - experience: 工作经验（从read_experience工具获取的内容）
   - skills: 技能信息（从read_skills工具获取的内容）
   - education: 教育背景（从read_education工具获取的内容）
   - extras: 其他信息（从read_extras工具获取的内容）
5. 给出具体的、可操作的建议

请始终显示你的完整思考过程，让用户了解你是如何得出结论的。当用户提供职位描述时，必须立即调工具获取候选人信息。分析. when asked to compose resume 必须调用compose_resume工具生成最终的简历报告。
call all the tools you need and are available to you. do not hesitate. call the tools, dont wait for confirmation。
when composing resume give complete jsons`),
	}

	fmt.Println("=== LLM 聊天机器人 (流式响应) ===")
	fmt.Println("输入 'quit' 或 'exit' 退出聊天")
	fmt.Println("================================")
	fmt.Println()

	// 聊天循环
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("用户: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" || input == "exit" {
			fmt.Println("助手: 再见！")
			break
		}

		// 添加用户消息到历史
		messages = append(messages, schema.UserMessage(input))

		// 使用agent调用并流式输出响应
		fmt.Print("助手: ")
		streamResp, err := agent.Stream(ctx, messages)
		if err != nil {
			log.Printf("Agent调用错误: %v", err)
			continue
		}

		// 处理流式响应
		var assistantContent strings.Builder
		for {
			chunk, err := streamResp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("流式响应错误: %v", err)
				break
			}

			// 安全地输出chunk内容
			if chunk != nil && chunk.Content != "" {
				fmt.Print(chunk.Content)
				assistantContent.WriteString(chunk.Content)
			}
		}
		fmt.Println() // 换行

		// 将助手的完整响应添加到消息历史
		messages = append(messages, schema.AssistantMessage(assistantContent.String(), nil))
		fmt.Println()
	}
}
