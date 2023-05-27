package slack

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

type SlackPayload struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	Channel   string `json:"channel"`
	IconEmoji string `json:"icon_emoji"`
}

type SlackClient struct {
	webhookURL string
	channel    string
	username   string
}

func NewSlackClient(webhookURL, channel, username string) *SlackClient {
	return &SlackClient{
		webhookURL: webhookURL,
		channel:    channel,
		username:   username,
	}
}

func (c *SlackClient) Push(headerText string, data map[string]interface{}, emoji string) error {
	payload := SlackPayload{
		Text:      headerText + "\n```" + mapToString(data) + "```",
		Username:  c.username,
		Channel:   c.channel,
		IconEmoji: emoji,
	}

	return c.sendPayload(payload)
}

func (c *SlackClient) PushText(text, emoji string) error {
	payload := SlackPayload{
		Text:      text,
		Username:  c.username,
		Channel:   c.channel,
		IconEmoji: emoji,
	}

	return c.sendPayload(payload)
}

func (c *SlackClient) sendPayload(payload SlackPayload) error {
	jsonPayload := []byte(`{"text":"` + payload.Text + `","username":"` + payload.Username + `","channel":"` + payload.Channel + `","icon_emoji":"` + payload.IconEmoji + `"}`)

	req, err := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func mapToString(data map[string]interface{}) string {
	var buffer bytes.Buffer
	for key, value := range data {
		buffer.WriteString("\n\t" + key + ": " + toString(value))
	}
	return buffer.String()
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
