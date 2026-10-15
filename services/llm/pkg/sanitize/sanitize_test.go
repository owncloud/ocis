package sanitize

import "testing"

var validMessages = []interface{}{map[string]interface{}{"role": "user", "content": "hi"}}

func TestBody_RequiresModel(t *testing.T) {
	_, err := Body(map[string]interface{}{
		"messages": validMessages,
	}, "", 4096)
	if err == nil {
		t.Fatal("expected an error when model is missing and no override is set")
	}
}

func TestBody_ModelOverrideWins(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":    "client-requested-model",
		"messages": validMessages,
	}, "forced-model", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "forced-model" {
		t.Fatalf("expected model override to win, got %q", req.Model)
	}
}

func TestBody_ClientModelUsedWhenNoOverride(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":    "client-model",
		"messages": validMessages,
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "client-model" {
		t.Fatalf("expected client model, got %q", req.Model)
	}
}

func TestBody_RequiresMessagesArray(t *testing.T) {
	_, err := Body(map[string]interface{}{
		"model":    "m",
		"messages": "not-an-array",
	}, "", 4096)
	if err == nil {
		t.Fatal("expected an error when messages is not an array")
	}
}

func TestBody_MaxTokensClampedToLimit(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":      "m",
		"messages":   validMessages,
		"max_tokens": float64(999999),
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.MaxTokens != 4096 {
		t.Fatalf("expected max_tokens clamped to 4096, got %d", req.MaxTokens)
	}
}

func TestBody_MaxTokensDefaultsToLimitWhenAbsent(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":    "m",
		"messages": validMessages,
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.MaxTokens != 4096 {
		t.Fatalf("expected max_tokens to default to the limit, got %d", req.MaxTokens)
	}
}

func TestBody_MaxTokensDefaultsToLimitWhenNotPositive(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":      "m",
		"messages":   validMessages,
		"max_tokens": float64(-1),
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.MaxTokens != 4096 {
		t.Fatalf("expected max_tokens to default to the limit for a non-positive value, got %d", req.MaxTokens)
	}
}

func TestBody_TemperaturePassedThrough(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":       "m",
		"messages":    validMessages,
		"temperature": float64(0.7),
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Temperature == nil || *req.Temperature != 0.7 {
		t.Fatalf("expected temperature 0.7, got %v", req.Temperature)
	}
}

func TestBody_TemperatureOmittedWhenAbsent(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":    "m",
		"messages": validMessages,
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Temperature != nil {
		t.Fatalf("expected temperature to be nil, got %v", req.Temperature)
	}
}

func TestBody_DropsUnknownFields(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":             "m",
		"messages":          validMessages,
		"some_random_field": "should be dropped",
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "m" || len(req.Messages) != 1 {
		t.Fatalf("unexpected sanitized request: %+v", req)
	}
}

func TestBody_RejectsInvalidMessages(t *testing.T) {
	tests := map[string][]interface{}{
		"empty":           {},
		"non-object":      {"hi"},
		"missing role":    {map[string]interface{}{"content": "hi"}},
		"non-string role": {map[string]interface{}{"role": 1}},
	}
	for name, messages := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Body(map[string]interface{}{"model": "m", "messages": messages}, "", 4096)
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestBody_RejectsNonPositiveMaxTokensLimit(t *testing.T) {
	for _, limit := range []int{0, -1} {
		_, err := Body(map[string]interface{}{"model": "m", "messages": validMessages}, "", limit)
		if err == nil {
			t.Fatalf("expected an error for limit %d", limit)
		}
	}
}

func TestBody_RejectsStreaming(t *testing.T) {
	_, err := Body(map[string]interface{}{"model": "m", "messages": validMessages, "stream": true}, "", 4096)
	if err == nil {
		t.Fatal("expected an error when stream is true")
	}
}

func TestBody_AllowsStreamFalse(t *testing.T) {
	_, err := Body(map[string]interface{}{"model": "m", "messages": validMessages, "stream": false}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
