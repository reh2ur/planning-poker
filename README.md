# Planning Poker

A free, ad-free, keyboard-first planning poker tool. Join from the browser or over SSH. Stores nothing.

Status: early development. See [docs/SPEC.md](docs/SPEC.md) for the v1 spec and the GitHub milestones for progress.

## Develop

Requirements: Go (version in `go.mod`), Node (see `.nvmrc`), pnpm, golangci-lint, lefthook, gitleaks.

```
cp .env.example .env
lefthook install
make dev-server   # terminal 1
make dev-web      # terminal 2
make lint test
```

Dev loop: run `make dev-server` and `make dev-web` in two terminals and open the Vite URL (http://localhost:5173). Vite hot-reloads the web app and proxies `/ws` and `/healthz` to the Go server. The proxy target is derived from `POKER_HTTP_ADDR` (default `localhost:8080`), so export it in both terminals to use another port.

## Run with Docker

```
cp .env.example .env
docker compose -f deploy/docker-compose.yml up --build
```

Serves the web app on http://localhost:8080. Configuration comes from `.env`; the published port in `deploy/docker-compose.yml` assumes the default `POKER_HTTP_ADDR`. The image contains no env or key files.

## License

MIT
