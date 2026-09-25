// Package sanitize validates and sanitizes the JSON body a client sends to
// the chat-completions endpoint before it is forwarded to the upstream LLM.
// Only the fields the LLM actually needs are forwarded; everything else the
// client sent is dropped, so a client cannot inject unexpected upstream
// parameters.
package sanitize

import "errors"

// Request is the sanitized body forwarded to the upstream LLM.
type Request struct {
	Model       string        `json:"model"`
	Messages    []interface{} `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature *float64      `json:"temperature,omitempty"`
}

// Body validates and sanitizes a client's decoded JSON body.
//
// modelOverride, when non-empty, replaces whatever model the client sent.
// maxTokensLimit is a hard ceiling: the client's max_tokens is used only if
// positive, and is always clamped to this limit.
func Body(raw map[string]interface{}, modelOverride string, maxTokensLimit int) (Request, error) {
	model := modelOverride
	if model == "" {
		if m, ok := raw["model"].(string); ok {
			model = m
		}
	}
	if model == "" {
		return Request{}, errors.New("model is required")
	}

	messages, ok := raw["messages"].([]interface{})
	if !ok {
		return Request{}, errors.New("messages must be an array")
	}

	maxTokens := maxTokensLimit
	if v, ok := raw["max_tokens"].(float64); ok && v > 0 {
		maxTokens = int(v)
	}
	if maxTokens > maxTokensLimit {
		maxTokens = maxTokensLimit
	}

	req := Request{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}

	if t, ok := raw["temperature"].(float64); ok {
		req.Temperature = &t
	}

	return req, nil
}
