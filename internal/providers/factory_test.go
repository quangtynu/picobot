package providers

import (
	"testing"

	"github.com/local/picobot/internal/config"
)

func TestNewProviderFromConfig_PicksOllama(t *testing.T) {
	cfg := config.Config{}
	cfg.Providers.Ollama = &config.ProviderConfig{APIBase: "http://localhost:11434/v1"}
	p := NewProviderFromConfig(cfg)
	_, ok := p.(*OllamaProvider)
	if !ok {
		t.Fatalf("expected OllamaProvider, got %T", p)
	}
}

func TestNewProviderFromConfig_PicksOpenRouter(t *testing.T) {
	cfg := config.Config{}
	cfg.Providers.OpenRouter = &config.ProviderConfig{APIKey: "test"}
	p := NewProviderFromConfig(cfg)
	_, ok := p.(*OpenRouterProvider)
	if !ok {
		t.Fatalf("expected OpenRouterProvider, got %T", p)
	}
}

func TestNewProviderFromConfig_FallbacksToStub(t *testing.T) {
	cfg := config.Config{}
	p := NewProviderFromConfig(cfg)
	_, ok := p.(*StubProvider)
	if !ok {
		t.Fatalf("expected StubProvider, got %T", p)
	}
}
