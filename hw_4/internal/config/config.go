package config

import "os"

type Config struct {
	Email      string
	PushDevice string
}

func Load() Config {
	return Config{
		Email:      env("NOTIFIER_EMAIL", "admin@example.com"),
		PushDevice: env("NOTIFIER_PUSH_DEVICE", "device-1"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
