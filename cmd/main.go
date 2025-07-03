package main

import (
	"context"
	"log"

	"wecom-bot-server-go/internal/config"
	"wecom-bot-server-go/internal/server"
	"wecom-bot-server-go/internal/wecom"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建企业微信机器人客户端
	wecomClient := wecom.NewClient(cfg.WeComWebhookKey)

	// 创建 MCP 服务器
	mcpServer := mcpserver.NewMCPServer(
		"wecom-bot-server",
		"2.0.0",
		mcpserver.WithToolCapabilities(false),
		mcpserver.WithRecovery(),
		mcpserver.WithLogging(),
	)

	// 创建服务器实例并注册工具
	srv := server.New(mcpServer, wecomClient)
	if err := srv.RegisterTools(context.Background()); err != nil {
		log.Fatalf("注册工具失败: %v", err)
	}

	// 启动服务器
	log.Println("启动企业微信机器人 MCP 服务器...")
	if err := mcpserver.ServeStdio(mcpServer); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}
