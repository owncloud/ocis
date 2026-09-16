package parser

import (
	"testing"

	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
)

func TestValidate_RequiresEndpoint(t *testing.T) {
	cfg := &config.Config{}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected an error when LLM.Endpoint is empty")
	}
}

func TestValidate_PassesWhenEndpointSet(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLM{
			Endpoint: "https://api.example.com/v1",
		},
	}

	if err := Validate(cfg); err != nil {
		t.Fatalf("expected no error when LLM.Endpoint is set, got: %v", err)
	}
}
