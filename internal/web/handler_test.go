package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const shell = "<!doctype html><title>shell</title>"

func builtFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html":                    {Data: []byte(shell)},
		"robots.txt":                    {Data: []byte("User-agent: *")},
		"_app/immutable/entry/app.js":   {Data: []byte("export {};")},
		"_app/immutable/assets/app.css": {Data: []byte("body{}")},
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name     string
		fs       fstest.MapFS
		method   string
		path     string
		wantCode int
		wantBody string
	}{
		{"root serves shell", builtFS(), http.MethodGet, "/", http.StatusOK, shell},
		{"room url serves shell", builtFS(), http.MethodGet, "/abc123", http.StatusOK, shell},
		{"nested room url serves shell", builtFS(), http.MethodGet, "/r/abc123", http.StatusOK, shell},
		{"existing asset is served", builtFS(), http.MethodGet, "/_app/immutable/entry/app.js", http.StatusOK, "export {};"},
		{"existing root file is served", builtFS(), http.MethodGet, "/robots.txt", http.StatusOK, "User-agent: *"},
		{"missing app asset is 404", builtFS(), http.MethodGet, "/_app/missing.js", http.StatusNotFound, ""},
		{"missing app path without extension is 404", builtFS(), http.MethodGet, "/_app/immutable/missing", http.StatusNotFound, ""},
		{"missing image is 404", builtFS(), http.MethodGet, "/x.png", http.StatusNotFound, ""},
		{"head on room url", builtFS(), http.MethodHead, "/abc123", http.StatusOK, ""},
		{"post is rejected", builtFS(), http.MethodPost, "/abc123", http.StatusMethodNotAllowed, ""},
		{"healthz still works", builtFS(), http.MethodGet, "/healthz", http.StatusOK, "ok"},
		{"healthz works without build", fstest.MapFS{}, http.MethodGet, "/healthz", http.StatusOK, "ok"},
		{"root without build is 503", fstest.MapFS{}, http.MethodGet, "/", http.StatusServiceUnavailable, "web UI not built"},
		{"room url without build is 503", fstest.MapFS{}, http.MethodGet, "/abc123", http.StatusServiceUnavailable, "web UI not built"},
		{"asset without build is 404", fstest.MapFS{}, http.MethodGet, "/x.png", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			NewHandler(tt.fs).ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d (body %q)", rec.Code, tt.wantCode, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), tt.wantBody)
			}
			if tt.wantCode == http.StatusNotFound && strings.Contains(rec.Body.String(), "shell") {
				t.Fatalf("404 must not return the shell: %q", rec.Body.String())
			}
		})
	}
}

func TestShellContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHandler(builtFS()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/abc123", nil))
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("content-type = %q, want text/html", got)
	}
}

func TestEmbeddedHandlerServesHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("got %d %q", rec.Code, rec.Body.String())
	}
}
