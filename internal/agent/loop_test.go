package agent

import (
	"testing"
	"time"

	"github.com/local/picobot/internal/bus"
	"github.com/local/picobot/internal/providers"
)

func TestProcessDirectWithStub(t *testing.T) {
	b := bus.NewMessageBus(10)
	p := providers.NewStubProvider()

	ag := NewAgentLoop(b, p, p.GetDefaultModel(), 5, "", nil)

	resp, err := ag.ProcessDirect("hello", 1*time.Second)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == "" {
		t.Fatalf("expected response, got empty string")
	}
}
