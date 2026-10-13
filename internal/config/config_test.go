package config

import (
	"testing"
	"time"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
		check   func(*testing.T, Config)
	}{
		{
			name: "defaults",
			env:  map[string]string{},
			check: func(t *testing.T, c Config) {
				if c.Env != "development" || c.HTTPAddr != ":8080" || c.AwayTimeout != 3*time.Minute {
					t.Fatalf("unexpected defaults: %+v", c)
				}
			},
		},
		{
			name:    "production without host key",
			env:     map[string]string{"POKER_ENV": "production"},
			wantErr: true,
		},
		{
			name: "production with host key",
			env:  map[string]string{"POKER_ENV": "production", "POKER_SSH_HOST_KEY_PATH": "/keys/host"},
			check: func(t *testing.T, c Config) {
				if !c.Production() || c.SSHHostKeyPath != "/keys/host" {
					t.Fatalf("unexpected config: %+v", c)
				}
			},
		},
		{
			name:    "bad away timeout",
			env:     map[string]string{"POKER_AWAY_TIMEOUT": "soon"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Load(env(tt.env))
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}
