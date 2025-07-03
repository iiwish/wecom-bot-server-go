# 企业微信机器人 MCP 服务器

这是一个基于 Go 语言的 MCP (Model Context Protocol) 服务器，提供企业微信机器人功能。

## 功能特性

- 🔧 **发送文本消息** - 支持 @ 用户功能
- 📝 **发送 Markdown 消息** - 支持富文本格式
- 🖼️ **发送图片消息** - 支持 Base64 编码的图片
- 📰 **发送图文消息** - 支持链接预览
- 🎴 **发送模板卡片** - 支持交互式卡片
- 📁 **文件上传** - 支持各种文件格式

## 安装和配置

### 1. 克隆项目

```bash
git clone https://github.com/iiwish/wecom-bot-server-go
cd wecom-bot-server-go
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置方式

无需配置文件或环境变量，所有企业微信 webhook_key 均通过每次工具调用参数传递，支持多机器人灵活适配。

### 4. 构建项目

```bash
go build ./cmd/main.go
```

### 5. 运行服务器

```bash
go run ./cmd/main.go
```

服务器启动后，将以 HTTP 协议监听 8080 端口，MCP 客户端可通过如下方式调用：

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"tool":"send_text","arguments":{"webhook_key":"xxx","content":"Hello"}}'
```

## MCP 工具说明

### send_text
发送文本消息到企业微信群

**参数：**
- `content` (必需): 要发送的文本内容
- `mentioned_list` (可选): 要@的用户ID列表，多个用户用逗号分隔
- `mentioned_mobile_list` (可选): 要@的手机号列表，多个用户用逗号分隔

**示例：**
```json
{
  "content": "Hello, 大家好！",
  "mentioned_list": "@xiaoyang,@wike",
  "mentioned_mobile_list": "13812345678,13987654321"
}
```

### send_markdown
发送 Markdown 消息到企业微信群

**参数：**
- `content` (必需): 要发送的 Markdown 内容

**示例：**
```json
{
  "content": "## 项目更新\n\n- [x] 完成功能A\n- [ ] 开发功能B\n\n**详情：** [查看链接](https://example.com)"
}
```

### send_image
发送图片消息到企业微信群

**参数：**
- `base64_data` (必需): Base64编码的图片数据
- `md5` (必需): 图片的MD5哈希值

### send_news
发送图文消息到企业微信群

**参数：**
- `title` (必需): 图文消息标题
- `description` (可选): 图文消息描述
- `url` (必需): 图文消息链接
- `picurl` (可选): 图文消息图片链接

**示例：**
```json
{
  "title": "新产品发布",
  "description": "我们很高兴宣布新产品正式发布",
  "url": "https://example.com/product",
  "picurl": "https://example.com/image.jpg"
}
```

### send_template_card
发送模板卡片消息到企业微信群

**参数：**
- `card_type` (必需): 模板卡片类型
- `main_title` (必需): 模板卡片主标题
- `main_desc` (可选): 模板卡片主描述
- `card_action_type` (必需): 卡片动作类型
- `card_action_url` (可选): 卡片动作URL
- `card_action_appid` (可选): 卡片动作应用ID
- `card_action_pagepath` (可选): 卡片动作页面路径

### upload_file
上传文件到企业微信

**参数：**
- `file_path` (必需): 要上传的文件路径

## 项目结构

```
wecom-bot-server-go/
├── cmd/
│   └── main.go              # 程序入口
├── internal/
│   ├── server/
│   │   └── server.go        # MCP 服务器实现
│   └── wecom/
│       └── client.go        # 企业微信客户端
├── go.mod                   # Go 模块文件
└── README.md               # 项目说明
```

## 在 Claude Desktop 中使用

在 Claude Desktop 的配置文件中添加以下配置：

### Windows
文件位置：`%APPDATA%\Claude\claude_desktop_config.json`

### macOS
文件位置：`~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "wecom-bot": {
      "command": "path/to/wecom-bot-server",
      "args": [],
      "cwd": "path/to/wecom-bot-server-go"
    }
  }
}
```

## 获取企业微信 Webhook Key

1. 登录企业微信管理后台
2. 进入"应用管理" -> "群机器人"
3. 创建新的群机器人或编辑现有机器人
4. 复制 Webhook URL 中的 `key` 参数值

> 使用本MCP工具时，每次调用需传递 webhook_key 参数，无需任何本地配置文件。

Webhook URL 格式：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY_HERE`

## 许可证


## 贡献

欢迎提交 Issue 和 Pull Request！

## 支持

如有问题，请创建 Issue 或联系维护者。
## License

This project is licensed under the MIT License.
