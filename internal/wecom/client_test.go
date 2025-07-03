// 文件名: client_test.go
package wecom

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"
)

const (
	testWebhookKey = "ce3ba63d-36bb-4681-8a07-0426557fc25a"
)

func TestSendText(t *testing.T) {
	client := NewClient(testWebhookKey)
	err := client.SendText("测试文本消息", []string{"user1", "user2"}, []string{"13800000000", "13900000000"})
	if err != nil {
		t.Fatalf("SendText failed: %v", err)
	}
}

func TestSendMarkdown(t *testing.T) {
	client := NewClient(testWebhookKey)
	err := client.SendMarkdown("### 测试Markdown消息\n> 这是一条测试")
	if err != nil {
		t.Fatalf("SendMarkdown failed: %v", err)
	}
}

func TestSendImage(t *testing.T) {
	client := NewClient(testWebhookKey)
	filePath := "test.jpg"
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		t.Skip("test.jpg 文件不存在，跳过图片发送测试")
		return
	}
	if err != nil {
		t.Fatalf("读取 test.jpg 失败: %v", err)
	}
	base64Data := encodeToBase64(data)
	md5Str := calcMD5(data)
	err = client.SendImage(base64Data, md5Str)
	if err != nil {
		t.Fatalf("SendImage failed: %v", err)
	}
}

func encodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func calcMD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func TestSendNews(t *testing.T) {
	client := NewClient(testWebhookKey)
	articles := []NewsArticle{
		{
			Title:       "测试图文标题",
			Description: "图文描述",
			URL:         "https://example.com",
			PicURL:      "https://example.com/image.png",
		},
	}
	err := client.SendNews(articles)
	if err != nil {
		t.Fatalf("SendNews failed: %v", err)
	}
}

func TestSendTemplateCard(t *testing.T) {
	client := NewClient(testWebhookKey)
	params := TemplateCardParams{
		CardType:           "text_notice",
		MainTitle:          "模板卡片主标题",
		MainDesc:           "主描述内容",
		CardActionType:     1,
		CardActionURL:      "https://example.com",
		CardActionAppID:    "",
		CardActionPagePath: "",
	}
	err := client.SendTemplateCard(params)
	if err != nil {
		t.Fatalf("SendTemplateCard failed: %v", err)
	}
}

func TestUploadFile(t *testing.T) {
	client := NewClient(testWebhookKey)
	// 请确保 test.txt 文件存在于项目根目录
	_, err := os.Stat("test.txt")
	if os.IsNotExist(err) {
		t.Skip("test.txt 文件不存在，跳过上传文件测试")
	}
	mediaID, err := client.UploadFile("test.txt")
	if err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}
	if mediaID == "" {
		t.Fatalf("UploadFile 返回空 mediaID")
	}
}
