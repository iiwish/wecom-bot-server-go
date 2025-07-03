package server

import (
	"context"
	"strings"

	"wecom-bot-server-go/internal/wecom"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server MCP服务器包装器
type Server struct {
	mcpServer *server.MCPServer
}

// New 创建新的服务器实例
func New(mcpServer *server.MCPServer) *Server {
	return &Server{
		mcpServer: mcpServer,
	}
}

// RegisterTools 注册所有工具
func (s *Server) RegisterTools(ctx context.Context) error {
	// 注册发送文本消息工具
	if err := s.registerSendTextTool(); err != nil {
		return err
	}

	// 注册发送Markdown消息工具
	if err := s.registerSendMarkdownTool(); err != nil {
		return err
	}

	// 注册发送图片消息工具
	if err := s.registerSendImageTool(); err != nil {
		return err
	}

	// 注册发送图文消息工具
	if err := s.registerSendNewsTool(); err != nil {
		return err
	}

	// 注册发送模板卡片工具
	if err := s.registerSendTemplateCardTool(); err != nil {
		return err
	}

	// 注册上传文件工具
	if err := s.registerUploadFileTool(); err != nil {
		return err
	}

	return nil
}

// registerSendTextTool 注册发送文本消息工具
func (s *Server) registerSendTextTool() error {
	tool := mcp.NewTool("send_text",
		mcp.WithDescription("发送文本消息到企业微信群"),
		mcp.WithString("webhook_key",
			mcp.Required(),
			mcp.Description("企业微信机器人的Webhook Key"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("要发送的文本内容"),
		),
		mcp.WithString("mentioned_list",
			mcp.Description("要@的用户ID列表，多个用户用逗号分隔，例如：@xiaoyang,@wike"),
		),
		mcp.WithString("mentioned_mobile_list",
			mcp.Description("要@的手机号列表，多个用户用逗号分隔，例如：@xiaoyang,@wike"),
		),
	)

	s.mcpServer.AddTool(tool, s.handleSendText)
	return nil
}

// registerSendMarkdownTool 注册发送Markdown消息工具
func (s *Server) registerSendMarkdownTool() error {
	tool := mcp.NewTool("send_markdown",
		mcp.WithDescription("发送Markdown消息到企业微信群"),
		mcp.WithString("webhook_key",
			mcp.Required(),
			mcp.Description("企业微信机器人的Webhook Key"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("要发送的Markdown内容"),
		),
	)

	s.mcpServer.AddTool(tool, s.handleSendMarkdown)
	return nil
}

// registerSendImageTool 注册发送图片消息工具
func (s *Server) registerSendImageTool() error {
	tool := mcp.NewTool("send_image",
		mcp.WithDescription("发送图片消息到企业微信群"),
		mcp.WithString("webhook_key",
			mcp.Required(),
			mcp.Description("企业微信机器人的Webhook Key"),
		),
		mcp.WithString("base64_data",
			mcp.Required(),
			mcp.Description("Base64编码的图片数据"),
		),
		mcp.WithString("md5",
			mcp.Required(),
			mcp.Description("图片的MD5哈希值"),
		),
	)

	s.mcpServer.AddTool(tool, s.handleSendImage)
	return nil
}

// registerSendNewsTool 注册发送图文消息工具
func (s *Server) registerSendNewsTool() error {
	tool := mcp.NewTool("send_news",
		mcp.WithDescription("发送图文消息到企业微信群"),
		mcp.WithString("webhook_key",
			mcp.Required(),
			mcp.Description("企业微信机器人的Webhook Key"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("图文消息标题"),
		),
		mcp.WithString("description",
			mcp.Description("图文消息描述"),
		),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("图文消息链接"),
		),
		mcp.WithString("picurl",
			mcp.Description("图文消息图片链接"),
		),
	)

	s.mcpServer.AddTool(tool, s.handleSendNews)
	return nil
}

// registerSendTemplateCardTool 注册发送模板卡片工具
func (s *Server) registerSendTemplateCardTool() error {
	tool := mcp.NewTool("send_template_card",
		mcp.WithDescription("发送模板卡片消息到企业微信群"),
		mcp.WithString("webhook_key",
			mcp.Required(),
			mcp.Description("企业微信机器人的Webhook Key"),
		),
		mcp.WithString("card_type",
			mcp.Required(),
			mcp.Description("模板卡片类型"),
		),
		mcp.WithString("main_title",
			mcp.Required(),
			mcp.Description("模板卡片主标题"),
		),
		mcp.WithString("main_desc",
			mcp.Description("模板卡片主描述"),
		),
		mcp.WithNumber("card_action_type",
			mcp.Required(),
			mcp.Description("卡片动作类型"),
		),
		mcp.WithString("card_action_url",
			mcp.Description("卡片动作URL"),
		),
		mcp.WithString("card_action_appid",
			mcp.Description("卡片动作应用ID"),
		),
		mcp.WithString("card_action_pagepath",
			mcp.Description("卡片动作页面路径"),
		),
	)

	s.mcpServer.AddTool(tool, s.handleSendTemplateCard)
	return nil
}

// registerUploadFileTool 注册上传文件工具
func (s *Server) registerUploadFileTool() error {
	tool := mcp.NewTool("upload_file",
		mcp.WithDescription("上传文件到企业微信"),
		mcp.WithString("webhook_key",
			mcp.Required(),
			mcp.Description("企业微信机器人的Webhook Key"),
		),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("要上传的文件路径"),
		),
	)

	s.mcpServer.AddTool(tool, s.handleUploadFile)
	return nil
}

// handleSendText 处理发送文本消息
func (s *Server) handleSendText(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	webhookKey, ok := args["webhook_key"].(string)
	if !ok || webhookKey == "" {
		return mcp.NewToolResultError("webhook_key参数必须是非空字符串"), nil
	}

	content, ok := args["content"].(string)
	if !ok {
		return mcp.NewToolResultError("content参数必须是字符串"), nil
	}

	var mentionedList []string
	if mentionedListStr, ok := args["mentioned_list"].(string); ok && mentionedListStr != "" {
		mentionedList = strings.Split(mentionedListStr, ",")
		for i, mention := range mentionedList {
			mentionedList[i] = strings.TrimSpace(mention)
		}
	}

	var mentionedMobileList []string
	if mentionedMobileListStr, ok := args["mentioned_mobile_list"].(string); ok && mentionedMobileListStr != "" {
		mentionedMobileList = strings.Split(mentionedMobileListStr, ",")
		for i, mobile := range mentionedMobileList {
			mentionedMobileList[i] = strings.TrimSpace(mobile)
		}
	}

	// 动态创建wecom客户端
	wecomClient := wecom.NewClient(webhookKey)
	err := wecomClient.SendText(content, mentionedList, mentionedMobileList)
	if err != nil {
		return mcp.NewToolResultError("发送文本消息失败: " + err.Error()), nil
	}

	return mcp.NewToolResultText("文本消息发送成功"), nil
}

// handleSendMarkdown 处理发送Markdown消息
func (s *Server) handleSendMarkdown(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	webhookKey, ok := args["webhook_key"].(string)
	if !ok || webhookKey == "" {
		return mcp.NewToolResultError("webhook_key参数必须是非空字符串"), nil
	}

	content, ok := args["content"].(string)
	if !ok {
		return mcp.NewToolResultError("content参数必须是字符串"), nil
	}

	// 动态创建wecom客户端
	wecomClient := wecom.NewClient(webhookKey)
	err := wecomClient.SendMarkdown(content)
	if err != nil {
		return mcp.NewToolResultError("发送Markdown消息失败: " + err.Error()), nil
	}

	return mcp.NewToolResultText("Markdown消息发送成功"), nil
}

// handleSendImage 处理发送图片消息
func (s *Server) handleSendImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	webhookKey, ok := args["webhook_key"].(string)
	if !ok || webhookKey == "" {
		return mcp.NewToolResultError("webhook_key参数必须是非空字符串"), nil
	}

	base64Data, ok := args["base64_data"].(string)
	if !ok {
		return mcp.NewToolResultError("base64_data参数必须是字符串"), nil
	}

	md5Hash, ok := args["md5"].(string)
	if !ok {
		return mcp.NewToolResultError("md5参数必须是字符串"), nil
	}

	// 动态创建wecom客户端
	wecomClient := wecom.NewClient(webhookKey)
	err := wecomClient.SendImage(base64Data, md5Hash)
	if err != nil {
		return mcp.NewToolResultError("发送图片消息失败: " + err.Error()), nil
	}

	return mcp.NewToolResultText("图片消息发送成功"), nil
}

// handleSendNews 处理发送图文消息
func (s *Server) handleSendNews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	webhookKey, ok := args["webhook_key"].(string)
	if !ok || webhookKey == "" {
		return mcp.NewToolResultError("webhook_key参数必须是非空字符串"), nil
	}

	title, ok := args["title"].(string)
	if !ok {
		return mcp.NewToolResultError("title参数必须是字符串"), nil
	}

	url, ok := args["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url参数必须是字符串"), nil
	}

	description := ""
	if desc, ok := args["description"].(string); ok {
		description = desc
	}

	picurl := ""
	if pic, ok := args["picurl"].(string); ok {
		picurl = pic
	}

	// 创建文章数组
	articles := []wecom.NewsArticle{
		{
			Title:       title,
			Description: description,
			URL:         url,
			PicURL:      picurl,
		},
	}

	// 动态创建wecom客户端
	wecomClient := wecom.NewClient(webhookKey)
	err := wecomClient.SendNews(articles)
	if err != nil {
		return mcp.NewToolResultError("发送图文消息失败: " + err.Error()), nil
	}

	return mcp.NewToolResultText("图文消息发送成功"), nil
}

// handleSendTemplateCard 处理发送模板卡片消息
func (s *Server) handleSendTemplateCard(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	webhookKey, ok := args["webhook_key"].(string)
	if !ok || webhookKey == "" {
		return mcp.NewToolResultError("webhook_key参数必须是非空字符串"), nil
	}

	cardType, ok := args["card_type"].(string)
	if !ok {
		return mcp.NewToolResultError("card_type参数必须是字符串"), nil
	}

	mainTitle, ok := args["main_title"].(string)
	if !ok {
		return mcp.NewToolResultError("main_title参数必须是字符串"), nil
	}

	cardActionType, ok := args["card_action_type"].(float64)
	if !ok {
		return mcp.NewToolResultError("card_action_type参数必须是数字"), nil
	}

	mainDesc := ""
	if desc, ok := args["main_desc"].(string); ok {
		mainDesc = desc
	}

	cardActionURL := ""
	if url, ok := args["card_action_url"].(string); ok {
		cardActionURL = url
	}

	cardActionAppID := ""
	if appid, ok := args["card_action_appid"].(string); ok {
		cardActionAppID = appid
	}

	cardActionPagePath := ""
	if pagepath, ok := args["card_action_pagepath"].(string); ok {
		cardActionPagePath = pagepath
	}

	params := wecom.TemplateCardParams{
		CardType:           cardType,
		MainTitle:          mainTitle,
		MainDesc:           mainDesc,
		CardActionType:     int(cardActionType),
		CardActionURL:      cardActionURL,
		CardActionAppID:    cardActionAppID,
		CardActionPagePath: cardActionPagePath,
	}

	// 动态创建wecom客户端
	wecomClient := wecom.NewClient(webhookKey)
	err := wecomClient.SendTemplateCard(params)
	if err != nil {
		return mcp.NewToolResultError("发送模板卡片消息失败: " + err.Error()), nil
	}

	return mcp.NewToolResultText("模板卡片消息发送成功"), nil
}

// handleUploadFile 处理上传文件
func (s *Server) handleUploadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	webhookKey, ok := args["webhook_key"].(string)
	if !ok || webhookKey == "" {
		return mcp.NewToolResultError("webhook_key参数必须是非空字符串"), nil
	}

	filePath, ok := args["file_path"].(string)
	if !ok {
		return mcp.NewToolResultError("file_path参数必须是字符串"), nil
	}

	// 动态创建wecom客户端
	wecomClient := wecom.NewClient(webhookKey)
	mediaID, err := wecomClient.UploadFile(filePath)
	if err != nil {
		return mcp.NewToolResultError("上传文件失败: " + err.Error()), nil
	}

	return mcp.NewToolResultText("文件上传成功，媒体ID: " + mediaID), nil
}
