package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ReadProjectsTool 读取项目经验工具
type ReadProjectsTool struct{}

// Info 返回工具信息
func (r *ReadProjectsTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "read_projects",
		Desc: "读取候选人的有关项目经验信息，可以读取所有项目或特定项目",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"project_name": {
				Desc:     "特定项目名称，如果不指定则返回所有项目",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具
func (r *ReadProjectsTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 显示接收到的参数
	fmt.Printf("[Tool Debug] ReadProjectsTool called with arguments: %s\n", argumentsInJSON)

	// 解析参数
	var params map[string]interface{}
	err := json.Unmarshal([]byte(argumentsInJSON), &params)
	if err != nil {
		return "", err
	}

	basePath := "blocks/projects/"
	projectName, hasProject := params["project_name"].(string)

	// 如果指定了特定项目，只读取该项目
	if hasProject {
		filePath := filepath.Join(basePath, projectName+".md")
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("读取项目文件失败: %v", err)
		}
		fmt.Printf("[Tool Debug] ReadProjectsTool returning project content, size: %d bytes\n", len(content))
		return string(content), nil
	}

	// 如果没有指定特定项目，读取所有项目文件
	allContent := r.readAllProjectFiles(basePath)
	fmt.Printf("[Tool Debug] ReadProjectsTool returning all projects content, size: %d bytes\n", len(allContent))
	return allContent, nil
}

// readAllProjectFiles 读取目录下所有项目文件
func (r *ReadProjectsTool) readAllProjectFiles(dirPath string) string {
	var allContent string

	// 读取目录下的所有文件
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("[Tool Debug] ReadProjectsTool failed to read directory: %v\n", err)
		return ""
	}

	// 遍历所有.md文件
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(dirPath, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}
			allContent += string(content) + "\n\n"
		}
	}

	return allContent
}
