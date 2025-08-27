package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ComposeResumeTool 简历组合工具
type ComposeResumeTool struct{}

// Info 返回工具信息
func (r *ComposeResumeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "compose_resume",
		Desc: "根据相关候选人信息以及挑选后的项目生成Markdown格式简历",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"job_title": {
				Desc:     "职位名称",
				Type:     schema.String,
				Required: true,
			},
			"candidate_name": {
				Desc:     "候选人姓名",
				Type:     schema.String,
				Required: true,
			},
			"relevant_projects": {
				Desc:     "与职位描述相关的项目内容",
				Type:     schema.String,
				Required: true,
			},
			"experience": {
				Desc:     "候选人的工作经验",
				Type:     schema.String,
				Required: false,
			},
			"skills": {
				Desc:     "候选人的技能信息",
				Type:     schema.String,
				Required: false,
			},
			"education": {
				Desc:     "候选人的教育背景",
				Type:     schema.String,
				Required: false,
			},
			"extras": {
				Desc:     "候选人的其他信息",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具
func (r *ComposeResumeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 显示接收到的参数
	fmt.Printf("[Tool Debug] ComposeResumeTool called with arguments: %s\n", argumentsInJSON)

	// 解析参数
	var params map[string]interface{}
	err := json.Unmarshal([]byte(argumentsInJSON), &params)
	if err != nil {
		return "", err
	}

	// 获取必需参数
	job_title, hasJobTitle := params["job_title"].(string)
	candidate_name, hasCandidateName := params["candidate_name"].(string)
	relevantProjects, hasProjects := params["relevant_projects"].(string)

	if !hasJobTitle || !hasCandidateName || !hasProjects {
		return "", fmt.Errorf("缺少必需参数: job_title, candidate_name 或 relevant_projects")
	}

	// 生成简历
	resumeContent := r.generateResume(
		job_title,
		relevantProjects,
		params["experience"],
		params["skills"],
		params["education"],
		params["extras"],
	)

	// 写入文件
	filename := fmt.Sprintf("outputs/%s_%s.md", job_title, candidate_name)
	err = os.WriteFile(filename, []byte(resumeContent), 0644)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("[Tool Debug] ComposeResumeTool returning resume content, size: %d bytes\n", len(resumeContent))
	return resumeContent, nil
}

// generateResume 生成简历
func (r *ComposeResumeTool) generateResume(
	job_title string,
	relevantProjects string,
	experience interface{},
	skills interface{},
	education interface{},
	extras interface{},
) string {
	var resume strings.Builder

	// 添加标题
	resume.WriteString("# 候选人简历\n\n")

	// 添加教育背景
	if education != nil && education != "" {
		resume.WriteString("## 教育背景\n")
		resume.WriteString(fmt.Sprintf("%v", education))
		resume.WriteString("\n\n")
	}

	// 添加工作经验
	if experience != nil && experience != "" {
		resume.WriteString("## 工作经验\n")
		resume.WriteString(fmt.Sprintf("%v", experience))
		resume.WriteString("\n\n")
	}

	// 添加技能
	if skills != nil && skills != "" {
		resume.WriteString("## 技能\n")
		resume.WriteString(fmt.Sprintf("%v", skills))
		resume.WriteString("\n\n")
	}

	// 添加相关项目经验
	if relevantProjects != "" {
		resume.WriteString("## 相关项目经验\n")
		resume.WriteString(relevantProjects)
		resume.WriteString("\n\n")
	}

	// 添加其他信息
	if extras != nil && extras != "" {
		resume.WriteString("## 其他信息\n")
		resume.WriteString(fmt.Sprintf("%v", extras))
		resume.WriteString("\n\n")
	}

	return resume.String()
}
