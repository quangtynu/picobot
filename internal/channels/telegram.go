package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/local/picobot/internal/bus"
)

// StartTelegram is a convenience wrapper that uses the real polling implementation
// with the standard Telegram base URL.
func StartTelegram(ctx context.Context, mb *bus.MessageBus, token string) error {
	if token == "" {
		return fmt.Errorf("telegram token not provided")
	}
	base := "https://api.telegram.org/bot" + token
	return StartTelegramWithBase(ctx, mb, token, base)
}

// StartTelegramWithBase starts long-polling against the given base URL (e.g., https://api.telegram.org/bot<TOKEN> or a test server URL).
func StartTelegramWithBase(ctx context.Context, mb *bus.MessageBus, token, base string) error {
	if base == "" {
		return fmt.Errorf("base URL is required")
	}

	client := &http.Client{Timeout: 45 * time.Second}

	// inbound polling goroutine
	go func() {
		offset := int64(0)
		for {
			select {
			case <-ctx.Done():
				log.Println("telegram: stopping inbound polling")
				return
			default:
			}

			values := url.Values{}
			values.Set("offset", strconv.FormatInt(offset, 10))
			values.Set("timeout", "30")
			u := base + "/getUpdates"
			resp, err := client.PostForm(u, values)
			if err != nil {
				log.Printf("telegram getUpdates error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var gu struct {
				Ok     bool `json:"ok"`
				Result []struct {
					UpdateID int64 `json:"update_id"`
					Message  *struct {
						MessageID int64 `json:"message_id"`
						From      *struct {
							ID int64 `json:"id"`
						} `json:"from"`
						Chat struct {
							ID int64 `json:"id"`
						} `json:"chat"`
						Text string `json:"text"`
					} `json:"message"`
				} `json:"result"`
			}
			if err := json.Unmarshal(body, &gu); err != nil {
				log.Printf("telegram: invalid getUpdates response: %v", err)
				continue
			}
			for _, upd := range gu.Result {
				if upd.UpdateID >= offset {
					offset = upd.UpdateID + 1
				}
				if upd.Message == nil {
					continue
				}
				m := upd.Message
				fromID := ""
				if m.From != nil {
					fromID = strconv.FormatInt(m.From.ID, 10)
				}
				chatID := strconv.FormatInt(m.Chat.ID, 10)
				mb.Inbound <- bus.InboundMessage{
					Channel:   "telegram",
					SenderID:  fromID,
					ChatID:    chatID,
					Content:   m.Text,
					Timestamp: time.Now(),
				}
			}
		}
	}()

	// outbound sender goroutine
	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		for {
			select {
			case <-ctx.Done():
				log.Println("telegram: stopping outbound sender")
				return
			case out := <-mb.Outbound:
				if out.Channel != "telegram" {
					// ignore messages for other channels
					continue
				}
				u := base + "/sendMessage"
				v := url.Values{}
				v.Set("chat_id", out.ChatID)
				v.Set("text", out.Content)
				resp, err := client.PostForm(u, v)
				if err != nil {
					log.Printf("telegram sendMessage error: %v", err)
					continue
				}
				io.ReadAll(resp.Body)
				resp.Body.Close()
			}
		}
	}()

	return nil
}
