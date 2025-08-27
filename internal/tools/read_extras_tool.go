package tools

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ReadExtrasTool 读取额外信息工具
type ReadExtrasTool struct{}

// Info 返回工具信息
func (r *ReadExtrasTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "read_extras",
		Desc: "读取候选人的有关额外信息（如证书、奖项、开源贡献等）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"file_path": {
				Desc:     "额外信息文件路径，默认为blocks/extras.md",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具
func (r *ReadExtrasTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 显示接收到的参数
	fmt.Printf("[Tool Debug] ReadExtrasTool called with arguments: %s\n", argumentsInJSON)

	// 解析参数
	filePath := "blocks/extras.md"

	// 这里可以解析JSON参数来获取自定义文件路径
	// 为了简化，我们直接使用默认路径

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 如果文件不存在，返回空内容而不是错误
		fmt.Printf("[Tool Debug] ReadExtrasTool file not found, returning empty content\n")
		return "", nil
	}

	fmt.Printf("[Tool Debug] ReadExtrasTool returning extras content, size: %d bytes\n", len(content))
	return string(content), nil
}
