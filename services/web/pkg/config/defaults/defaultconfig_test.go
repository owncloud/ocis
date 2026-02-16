package defaults

import "testing"

func TestSanitize_RemovesEmbedWhenAllUnset(t *testing.T) {
	cfg := DefaultConfig()

	Sanitize(cfg)

	// Assert: When all Embed fields are unset, Sanitize() should set Embed to nil.
	if cfg.Web.Config.Options.Embed != nil {
		// Assert field values are correct (all falsy)
		if cfg.Web.Config.Options.Embed.Enabled != nil {
			t.Errorf("Enabled should be nil, got %v", *cfg.Web.Config.Options.Embed.Enabled)
		}
		if cfg.Web.Config.Options.Embed.Target != "" {
			t.Errorf("Target should be empty, got %q", cfg.Web.Config.Options.Embed.Target)
		}
		if cfg.Web.Config.Options.Embed.MessagesOrigin != "" {
			t.Errorf("MessagesOrigin should be empty, got %q", cfg.Web.Config.Options.Embed.MessagesOrigin)
		}
		if cfg.Web.Config.Options.Embed.DelegateAuthentication {
			t.Errorf("DelegateAuthentication should be false, got %v", cfg.Web.Config.Options.Embed.DelegateAuthentication)
		}
		if cfg.Web.Config.Options.Embed.DelegateAuthenticationOrigin != "" {
			t.Errorf("DelegateAuthenticationOrigin should be empty, got %q", cfg.Web.Config.Options.Embed.DelegateAuthenticationOrigin)
		}

		// Assert: Embed struct value
		t.Errorf("Embed should be nil when all fields are unset, got %+v", cfg.Web.Config.Options.Embed)
	}
}
