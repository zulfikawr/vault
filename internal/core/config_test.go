package core

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test Defaults
	cfg := LoadConfig()
	if cfg.Port != 8090 {
		t.Errorf("expected default port 8090, got %d", cfg.Port)
	}

	// Test Env Overrides
	_ = os.Setenv("VAULT_PORT", "9090")
	defer func() { _ = os.Unsetenv("VAULT_PORT") }()

	cfg = LoadConfig()
	if cfg.Port != 9090 {
		t.Errorf("expected overridden port 9090, got %d", cfg.Port)
	}
}
