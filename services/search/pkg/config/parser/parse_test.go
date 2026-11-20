package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/config/defaults"
)

func validConfig() *config.Config {
	cfg := defaults.DefaultConfig()
	cfg.TokenManager = &config.TokenManager{JWTSecret: "secret"}
	cfg.ServiceAccount = config.ServiceAccount{ServiceAccountID: "id", ServiceAccountSecret: "secret"}
	return cfg
}

func TestDefaultMaxConsecutiveFailures(t *testing.T) {
	assert.Equal(t, 5, defaults.DefaultConfig().Extractor.MaxConsecutiveFailures)
}

func TestValidateMaxConsecutiveFailures(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   int
		wantErr bool
	}{
		{name: "default", value: 5},
		{name: "never abort", value: 0},
		{name: "negative", value: -1, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.Extractor.MaxConsecutiveFailures = tc.value
			err := Validate(cfg)
			if tc.wantErr {
				assert.ErrorContains(t, err, "SEARCH_EXTRACTOR_MAX_CONSECUTIVE_FAILURES")
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateFailedRetryInterval(t *testing.T) {
	assert.Zero(t, defaults.DefaultConfig().Extractor.FailedRetryInterval, "off by default: every walk retries a failed file")

	cfg := validConfig()
	cfg.Extractor.FailedRetryInterval = 24 * time.Hour
	assert.NoError(t, Validate(cfg))

	cfg.Extractor.FailedRetryInterval = -time.Second
	assert.ErrorContains(t, Validate(cfg), "SEARCH_EXTRACTOR_FAILED_RETRY_INTERVAL")
}
