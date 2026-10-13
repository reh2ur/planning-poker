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

## License

MIT
