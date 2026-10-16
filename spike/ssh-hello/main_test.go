package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2/testsession"
	gossh "golang.org/x/crypto/ssh"
)

func TestModelHandlesResizeAndQuit(t *testing.T) {
	tests := []struct {
		name string
		msgs []tea.Msg
		want string
	}{
		{"initial size", []tea.Msg{tea.WindowSizeMsg{Width: 80, Height: 24}}, "window: 80x24 (resizes: 1)"},
		{"second resize wins", []tea.Msg{
			tea.WindowSizeMsg{Width: 80, Height: 24},
			tea.WindowSizeMsg{Width: 100, Height: 30},
		}, "window: 100x30 (resizes: 2)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m tea.Model = model{}
			for _, msg := range tt.msgs {
				m, _ = m.Update(msg)
			}
			if got := m.View().Content; !strings.Contains(got, tt.want) {
				t.Fatalf("view %q does not contain %q", got, tt.want)
			}
		})
	}

	_, cmd := model{}.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if cmd == nil {
		t.Fatal("q must quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("q must produce QuitMsg")
	}
}

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func waitFor(t *testing.T, out *syncBuffer, want string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(out.String(), want) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %q in output %q", want, out.String())
}

func TestServerShowsAppAndHandlesResize(t *testing.T) {
	srv, err := newServer("127.0.0.1:0", "")
	if err != nil {
		t.Fatal(err)
	}
	sess := testsession.New(t, srv, &gossh.ClientConfig{
		User: "alex",
		Auth: []gossh.AuthMethod{gossh.Password("x")},
	})
	out := &syncBuffer{}
	sess.Stdout = out
	stdin, err := sess.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.RequestPty("xterm-256color", 24, 80, gossh.TerminalModes{}); err != nil {
		t.Fatal(err)
	}
	if err := sess.Shell(); err != nil {
		t.Fatal(err)
	}

	waitFor(t, out, "Hello from planning poker")
	waitFor(t, out, "80x24")

	if err := sess.WindowChange(40, 120); err != nil {
		t.Fatal(err)
	}
	waitFor(t, out, "120x40")

	if _, err := io.WriteString(stdin, "q"); err != nil {
		t.Fatal(err)
	}
	if err := sess.Wait(); err != nil {
		t.Fatalf("session did not exit cleanly: %v", err)
	}
}

func TestEphemeralHostKeyIsUnique(t *testing.T) {
	a, err := ephemeralHostKey()
	if err != nil {
		t.Fatal(err)
	}
	b, err := ephemeralHostKey()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(a, b) {
		t.Fatal("ephemeral keys must differ")
	}
}

func TestServerGeneratesMissingHostKeyAtConfiguredPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "host_key")
	if _, err := newServer("127.0.0.1:0", path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("host key not written at configured path: %v", err)
	}
}
