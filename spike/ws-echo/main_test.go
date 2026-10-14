package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func newTestServer(t *testing.T, logs io.Writer) *httptest.Server {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(logs, nil))
	srv := httptest.NewServer(newHandler(logger))
	t.Cleanup(srv.Close)
	return srv
}

func dial(t *testing.T, srv *httptest.Server) *websocket.Conn {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })
	return conn
}

func TestEchoReturnsMessages(t *testing.T) {
	srv := newTestServer(t, io.Discard)
	conn := dial(t, srv)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, msg := range []string{"ping 1", "", "héllo"} {
		if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
			t.Fatalf("write %q: %v", msg, err)
		}
		typ, got, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read after %q: %v", msg, err)
		}
		if typ != websocket.MessageText || string(got) != msg {
			t.Errorf("echo of %q = %v %q", msg, typ, got)
		}
	}
}

func TestServesPage(t *testing.T) {
	srv := newTestServer(t, io.Discard)
	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	for _, want := range []string{"reconnects", "interval"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("page missing %q", want)
		}
	}
}

func TestLogsConnectionLifetimeWithoutAddress(t *testing.T) {
	var logs syncBuffer
	srv := newTestServer(t, &logs)
	conn := dial(t, srv)
	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Logf("close: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for !strings.Contains(logs.String(), "ws closed") {
		if time.Now().After(deadline) {
			t.Fatalf("no close log; got %q", logs.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
	out := logs.String()
	for _, want := range []string{"ws opened", "duration="} {
		if !strings.Contains(out, want) {
			t.Errorf("log missing %q: %s", want, out)
		}
	}
	if strings.Contains(out, "127.0.0.1") {
		t.Errorf("log contains client address: %s", out)
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
