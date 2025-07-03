package wecom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	// WeComBotBaseURL 企业微信机器人基础URL
	WeComBotBaseURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook"
	// WeComBotSendURL 发送消息URL模板
	WeComBotSendURL = WeComBotBaseURL + "/send?key="
	// WeComBotUploadURL 上传文件URL模板
	WeComBotUploadURL = WeComBotBaseURL + "/upload_media?key="
)

// Client 企业微信机器人客户端
type Client struct {
	webhookKey string
	httpClient *http.Client
}

// NewClient 创建新的企业微信机器人客户端
func NewClient(webhookKey string) *Client {
	return &Client{
		webhookKey: webhookKey,
		httpClient: &http.Client{},
	}
}

// SendText 发送文本消息
func (c *Client) SendText(content string, mentionedList, mentionedMobileList []string) error {
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content":               content,
			"mentioned_list":        mentionedList,
			"mentioned_mobile_list": mentionedMobileList,
		},
	}

	return c.sendRequest(WeComBotSendURL, payload)
}

// SendMarkdown 发送Markdown消息
func (c *Client) SendMarkdown(content string) error {
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}

	return c.sendRequest(WeComBotSendURL, payload)
}

// SendImage 发送图片消息
func (c *Client) SendImage(base64Data, md5 string) error {
	payload := map[string]interface{}{
		"msgtype": "image",
		"image": map[string]interface{}{
			"base64": base64Data,
			"md5":    md5,
		},
	}

	return c.sendRequest(WeComBotSendURL, payload)
}

// NewsArticle 新闻文章结构
type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl"`
}

// SendNews 发送图文消息
func (c *Client) SendNews(articles []NewsArticle) error {
	payload := map[string]interface{}{
		"msgtype": "news",
		"news": map[string]interface{}{
			"articles": articles,
		},
	}

	return c.sendRequest(WeComBotSendURL, payload)
}

// TemplateCardParams 模板卡片参数
type TemplateCardParams struct {
	CardType           string
	MainTitle          string
	MainDesc           string
	CardActionType     int
	CardActionURL      string
	CardActionAppID    string
	CardActionPagePath string
}

// SendTemplateCard 发送模板卡片消息
func (c *Client) SendTemplateCard(params TemplateCardParams) error {
	payload := map[string]interface{}{
		"msgtype": "template_card",
		"template_card": map[string]interface{}{
			"card_type": params.CardType,
			"main_title": map[string]interface{}{
				"title": params.MainTitle,
				"desc":  params.MainDesc,
			},
			"card_action": map[string]interface{}{
				"type":     params.CardActionType,
				"url":      params.CardActionURL,
				"appid":    params.CardActionAppID,
				"pagepath": params.CardActionPagePath,
			},
		},
	}

	return c.sendRequest(WeComBotSendURL, payload)
}

// UploadFile 上传文件并返回媒体ID
func (c *Client) UploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("media", filePath)
	if err != nil {
		return "", fmt.Errorf("创建表单文件失败: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("复制文件内容失败: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭写入器失败: %w", err)
	}

	url := fmt.Sprintf("%s%s&type=file", WeComBotUploadURL, c.webhookKey)
	resp, err := c.httpClient.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return "", fmt.Errorf("上传文件请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		errMsg, _ := result["errmsg"].(string)
		return "", fmt.Errorf("企业微信API错误: %s", errMsg)
	}

	mediaID, ok := result["media_id"].(string)
	if !ok {
		return "", fmt.Errorf("无法获取媒体ID")
	}

	return mediaID, nil
}

// sendRequest 发送请求到企业微信API
func (c *Client) sendRequest(urlTemplate string, payload map[string]interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	url := urlTemplate + c.webhookKey
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		errMsg, _ := result["errmsg"].(string)
		return fmt.Errorf("企业微信API错误: %s", errMsg)
	}

	return nil
}
