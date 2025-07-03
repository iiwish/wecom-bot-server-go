package main

import (
	"context"
	"log"

	"wecom-bot-server-go/internal/server"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建 MCP 服务器
	mcpServer := mcpserver.NewMCPServer(
		"wecom-bot-server",
		"2.0.0",
		mcpserver.WithToolCapabilities(true),
		mcpserver.WithRecovery(),
		mcpserver.WithLogging(),
		mcpserver.WithInstructions("该工具支持通过企业微信机器人向群聊发送文本、Markdown、图片、图文、模板卡片等多种类型的消息，并支持文件上传。每次调用可灵活指定 webhook_key，无需本地配置，适用于多机器人、多群场景。适合自动化推送通知、播报信息、群内互动等企业微信场景。"),
	)

	// 创建服务器实例并注册工具
	srv := server.New(mcpServer)
	if err := srv.RegisterTools(context.Background()); err != nil {
		log.Fatalf("注册工具失败: %v", err)
	}

	// 启动服务器
	log.Println("启动企业微信机器人 MCP Streamable-HTTP 服务器，监听端口 :20301 ...")
	if err := mcpserver.NewStreamableHTTPServer(mcpServer, mcpserver.WithStateLess(true)).Start(":20301"); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}
