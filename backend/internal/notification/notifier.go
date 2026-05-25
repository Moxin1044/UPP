package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
)

type Notifier interface {
	Type() string
	Send(title, message string) error
}

// Lark (飞书)
type LarkNotifier struct {
	WebhookURL string `json:"webhook_url"`
}

func (n *LarkNotifier) Type() string { return "lark" }

func (n *LarkNotifier) Send(title, message string) error {
	payload := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]string{"tag": "plain_text", "content": title},
			},
			"elements": []map[string]interface{}{
				{"tag": "div", "text": map[string]string{"tag": "plain_text", "content": message}},
			},
		},
	}
	return sendWebhook(n.WebhookURL, payload)
}

// DingTalk (钉钉)
type DingTalkNotifier struct {
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

func (n *DingTalkNotifier) Type() string { return "dingtalk" }

func (n *DingTalkNotifier) Send(title, message string) error {
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  fmt.Sprintf("### %s\n%s", title, message),
		},
	}
	return sendWebhook(n.WebhookURL, payload)
}

// WebHook
type WebHookNotifier struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

func (n *WebHookNotifier) Type() string { return "webhook" }

func (n *WebHookNotifier) Send(title, message string) error {
	payload := map[string]string{"title": title, "message": message}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(strings.ToUpper(n.Method), n.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.Headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}

// Email
type EmailNotifier struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort string `json:"smtp_port"`
	From     string `json:"from"`
	Password string `json:"password"`
	To       string `json:"to"`
}

func (n *EmailNotifier) Type() string { return "email" }

func (n *EmailNotifier) Send(title, message string) error {
	auth := smtp.PlainAuth("", n.From, n.Password, n.SMTPHost)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", n.From, n.To, title, message)
	addr := n.SMTPHost + ":" + n.SMTPPort
	return smtp.SendMail(addr, auth, n.From, []string{n.To}, []byte(msg))
}

func sendWebhook(url string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
