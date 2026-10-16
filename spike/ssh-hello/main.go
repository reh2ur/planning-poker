// Throwaway spike for issue #13: Wish server with a tiny Bubble Tea app.
// Delete or move to docs after the validation ADR.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	wishtea "charm.land/wish/v2/bubbletea"

	"github.com/reh2ur/planning-poker/internal/config"
)

func main() {
	cfg, err := config.FromOS()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	srv, err := newServer(cfg.SSHAddr, cfg.SSHHostKeyPath)
	if err != nil {
		slog.Error("cannot create ssh server", "error", err)
		os.Exit(1)
	}
	slog.Info("listening", "addr", cfg.SSHAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// newServer reads the host key from hostKeyPath; an empty path means an
// ephemeral key, which is only acceptable for local runs.
func newServer(addr, hostKeyPath string) (*ssh.Server, error) {
	var keyOpt ssh.Option
	if hostKeyPath != "" {
		keyOpt = wish.WithHostKeyPath(hostKeyPath)
	} else {
		slog.Warn("POKER_SSH_HOST_KEY_PATH not set, using an ephemeral host key")
		pemBytes, err := ephemeralHostKey()
		if err != nil {
			return nil, err
		}
		keyOpt = wish.WithHostKeyPEM(pemBytes)
	}
	return wish.NewServer(
		wish.WithAddress(addr),
		keyOpt,
		wish.WithMiddleware(wishtea.Middleware(func(ssh.Session) (tea.Model, []tea.ProgramOption) {
			return model{}, nil
		})),
	)
}

func ephemeralHostKey() ([]byte, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate host key: %w", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("marshal host key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

type model struct {
	width, height int
	resizes       int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resizes++
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView(fmt.Sprintf(
		"Hello from planning poker\n\nwindow: %dx%d (resizes: %d)\n\npress q to quit\n",
		m.width, m.height, m.resizes))
	v.AltScreen = true
	return v
}
