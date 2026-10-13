package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// dist holds the SvelteKit build copied in by `make build-web`. Only the
// placeholder is tracked, so the package compiles without a web build.
//
//go:embed all:dist
var dist embed.FS

const shellFile = "index.html"

func Handler() http.Handler {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return NewHandler(sub)
}

// NewHandler serves static files from fsys. Unknown paths without a file
// extension (room URLs) get the app shell; unknown asset paths get a 404.
func NewHandler(fsys fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("GET /", static{fsys: fsys, files: http.FileServerFS(fsys)})
	return mux
}

type static struct {
	fsys  fs.FS
	files http.Handler
}

func (s static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if name == "" {
		s.serveShell(w, r)
		return
	}
	if fs.ValidPath(name) {
		if info, err := fs.Stat(s.fsys, name); err == nil && !info.IsDir() {
			s.files.ServeHTTP(w, r)
			return
		}
	}
	if isAssetPath(name) {
		http.NotFound(w, r)
		return
	}
	s.serveShell(w, r)
}

func (s static) serveShell(w http.ResponseWriter, r *http.Request) {
	if _, err := fs.Stat(s.fsys, shellFile); err != nil {
		http.Error(w, "web UI not built; run `make build`", http.StatusServiceUnavailable)
		return
	}
	http.ServeFileFS(w, r, s.fsys, shellFile)
}

func isAssetPath(name string) bool {
	return strings.HasPrefix(name, "_app/") || path.Ext(name) != ""
}
