package sanitize

import "testing"

func TestBody_RequiresModel(t *testing.T) {
	_, err := Body(map[string]interface{}{
		"messages": []interface{}{"hi"},
	}, "", 4096)
	if err == nil {
		t.Fatal("expected an error when model is missing and no override is set")
	}
}

func TestBody_ModelOverrideWins(t *testing.T) {
	req, err := Body(map[string]interface{}{
		"model":    "client-requested-model",
		"messages": []interface{}{"hi"},
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
		"messages": []interface{}{"hi"},
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
		"messages":   []interface{}{"hi"},
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
		"messages": []interface{}{"hi"},
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
		"messages":   []interface{}{"hi"},
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
		"messages":    []interface{}{"hi"},
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
		"messages": []interface{}{"hi"},
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
		"messages":          []interface{}{"hi"},
		"some_random_field": "should be dropped",
	}, "", 4096)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "m" || len(req.Messages) != 1 {
		t.Fatalf("unexpected sanitized request: %+v", req)
	}
}
