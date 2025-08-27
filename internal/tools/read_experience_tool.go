package tools

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ReadExperienceTool 读取工作经验工具
type ReadExperienceTool struct{}

// Info 返回工具信息
func (r *ReadExperienceTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "read_experience",
		Desc: "读取候选人的有关工作经验信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"file_path": {
				Desc:     "工作经验文件路径，默认为blocks/experience.md",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具
func (r *ReadExperienceTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 显示接收到的参数
	fmt.Printf("[Tool Debug] ReadExperienceTool called with arguments: %s\n", argumentsInJSON)

	// 解析参数
	filePath := "blocks/experience.md"

	// 这里可以解析JSON参数来获取自定义文件路径
	// 为了简化，我们直接使用默认路径

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	fmt.Printf("[Tool Debug] ReadExperienceTool returning experience content, size: %d bytes\n", len(content))
	return string(content), nil
}
