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
