package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sudhucodes/pingsystem/internal/sysinfo"
)

// EventType represents the event name.
type EventType string

const (
	EventStartup  EventType = "Startup / User Login"
	EventSleep    EventType = "Sleep"
	EventWake     EventType = "Wake"
	EventLock     EventType = "Lock"
	EventUnlock   EventType = "Unlock"
	EventShutdown EventType = "Shutdown"
)

// Client handles Telegram notifications.
type Client struct {
	botToken   string
	chatID     string
	httpClient *http.Client
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type sendMessageResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

// NewClient initializes a Telegram client.
func NewClient(botToken, chatID string) *Client {
	return &Client{
		botToken: botToken,
		chatID:   chatID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendEvent sends a system event notification to Telegram.
func (c *Client) SendEvent(eventType EventType, info sysinfo.Info) error {
	if c.botToken == "" || c.chatID == "" {
		return fmt.Errorf("telegram botToken or chatId is missing in configuration")
	}

	icon := getEventIcon(eventType)

	text := fmt.Sprintf(
		"%s *PingSystem Event: %s*\n\n"+
			"*Device Alias:* %s\n"+
			"*Username:* %s\n"+
			"*Timestamp:* %s",
		icon,
		escapeMarkdown(string(eventType)),
		escapeMarkdown(info.DeviceAlias),
		escapeMarkdown(info.Username),
		escapeMarkdown(info.Timestamp),
	)

	reqBody := sendMessageRequest{
		ChatID:    c.chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var tgResp sendMessageResponse
	if err := json.Unmarshal(body, &tgResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w (raw response: %s)", err, string(body))
	}

	if !tgResp.OK {
		return fmt.Errorf("telegram API error: %s", tgResp.Description)
	}

	return nil
}

func getEventIcon(eventType EventType) string {
	switch eventType {
	case EventStartup:
		return "🚀"
	case EventSleep:
		return "🌙"
	case EventWake:
		return "☀️"
	case EventLock:
		return "🔒"
	case EventUnlock:
		return "🔓"
	case EventShutdown:
		return "🛑"
	default:
		return "🔔"
	}
}

// escapeMarkdown escapes MarkdownV2 reserved characters in dynamic text strings.
func escapeMarkdown(text string) string {
	reserved := `\_*[]()~` + "`" + `>#+-=|{}.!`
	var b bytes.Buffer
	for _, r := range text {
		found := false
		for _, res := range reserved {
			if r == res {
				b.WriteRune('\\')
				b.WriteRune(r)
				found = true
				break
			}
		}
		if !found {
			b.WriteRune(r)
		}
	}
	return b.String()
}
