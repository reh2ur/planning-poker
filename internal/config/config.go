package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Env            string
	Domain         string
	HTTPAddr       string
	SSHAddr        string
	SSHHostKeyPath string
	AwayTimeout    time.Duration
}

func (c Config) Production() bool { return c.Env == "production" }

func Load(getenv func(string) string) (Config, error) {
	c := Config{
		Env:            valueOr(getenv("POKER_ENV"), "development"),
		Domain:         valueOr(getenv("POKER_DOMAIN"), "localhost"),
		HTTPAddr:       valueOr(getenv("POKER_HTTP_ADDR"), ":8080"),
		SSHAddr:        valueOr(getenv("POKER_SSH_ADDR"), ":2222"),
		SSHHostKeyPath: getenv("POKER_SSH_HOST_KEY_PATH"),
	}
	d, err := time.ParseDuration(valueOr(getenv("POKER_AWAY_TIMEOUT"), "3m"))
	if err != nil {
		return Config{}, errors.New("POKER_AWAY_TIMEOUT: invalid duration")
	}
	c.AwayTimeout = d
	if c.Production() && c.SSHHostKeyPath == "" {
		return Config{}, errors.New("POKER_SSH_HOST_KEY_PATH is required when POKER_ENV=production")
	}
	return c, nil
}

func FromOS() (Config, error) { return Load(os.Getenv) }

func valueOr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
