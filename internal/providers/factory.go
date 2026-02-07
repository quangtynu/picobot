package providers

import "github.com/local/picobot/internal/config"

// NewProviderFromConfig creates a provider based on the configuration.
// Simple rules (v0):
//   - if OpenRouter API key present -> OpenRouter
//   - else if Ollama APIBase present -> Ollama
//   - else fallback to stub
func NewProviderFromConfig(cfg config.Config) LLMProvider {
	if cfg.Providers.OpenRouter != nil && cfg.Providers.OpenRouter.APIKey != "" {
		return NewOpenRouterProvider(cfg.Providers.OpenRouter.APIKey, cfg.Providers.OpenRouter.APIBase)
	}
	if cfg.Providers.Ollama != nil && cfg.Providers.Ollama.APIBase != "" {
		return NewOllamaProvider(cfg.Providers.Ollama.APIBase)
	}
	return NewStubProvider()
}
