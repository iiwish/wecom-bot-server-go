package config

import (
	"fmt"
	"os"
)

// Config 包含应用程序配置
type Config struct {
	WeComWebhookKey string
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	webhookKey := os.Getenv("WECOM_BOT_WEBHOOK_KEY")
	if webhookKey == "" {
		return nil, fmt.Errorf("WECOM_BOT_WEBHOOK_KEY 环境变量是必需的")
	}

	return &Config{
		WeComWebhookKey: webhookKey,
	}, nil
}
