package tools

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ReadEducationTool 读取教育背景工具
type ReadEducationTool struct{}

// Info 返回工具信息
func (r *ReadEducationTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "read_education",
		Desc: "读取候选人的有关教育背景信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"file_path": {
				Desc:     "教育背景文件路径，默认为blocks/education.md",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具
func (r *ReadEducationTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 显示接收到的参数
	fmt.Printf("[Tool Debug] ReadEducationTool called with arguments: %s\n", argumentsInJSON)

	// 解析参数
	filePath := "blocks/education.md"

	// 这里可以解析JSON参数来获取自定义文件路径
	// 为了简化，我们直接使用默认路径

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	fmt.Printf("[Tool Debug] ReadEducationTool returning education content, size: %d bytes\n", len(content))
	return string(content), nil
}
