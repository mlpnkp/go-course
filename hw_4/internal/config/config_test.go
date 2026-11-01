package config_test

import (
	"testing"

	"go-course/hw_4/internal/config"
)

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("NOTIFIER_EMAIL", "anna@example.com")
	t.Setenv("NOTIFIER_PUSH_DEVICE", "device-7")

	cfg := config.Load()

	if cfg.Email != "anna@example.com" || cfg.PushDevice != "device-7" {
		t.Errorf("Load() = %+v, want values from environment", cfg)
	}
}

func TestLoadUsesDefaultsForEmptyEnvironment(t *testing.T) {
	t.Setenv("NOTIFIER_EMAIL", "")
	t.Setenv("NOTIFIER_PUSH_DEVICE", "")

	cfg := config.Load()

	if cfg.Email == "" || cfg.PushDevice == "" {
		t.Errorf("Load() = %+v, want defaults", cfg)
	}
}
