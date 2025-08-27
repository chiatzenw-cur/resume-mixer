package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ReadSkillsTool 读取技能工具
type ReadSkillsTool struct{}

// Info 返回工具信息
func (r *ReadSkillsTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "read_skills",
		Desc: "读取候选人的有关技能信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"file_path": {
				Desc:     "技能文件路径，默认为blocks/skills.md",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具
func (r *ReadSkillsTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 显示接收到的参数
	fmt.Printf("[Tool Debug] ReadSkillsTool called with arguments: %s\n", argumentsInJSON)

	// 解析参数
	filePath := "blocks/skills.md"

	// 这里可以解析JSON参数来获取自定义文件路径
	// 为了简化，我们直接使用默认路径

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	fmt.Printf("[Tool Debug] ReadSkillsTool returning skills content, size: %d bytes\n", len(content))
	return string(content), nil
}
