package providers

import (
	"context"
	"testing"
	"time"
)

func TestStubProviderEcho(t *testing.T) {
	p := NewStubProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	msgs := []Message{{Role: "user", Content: "hello world"}}
	resp, err := p.Chat(ctx, msgs, nil, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Content == "" {
		t.Fatalf("expected non-empty content")
	}
}
