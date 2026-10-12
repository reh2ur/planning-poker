# CLAUDE.md

Planning poker: a free, keyboard-first tool for the browser and SSH. Go server, SvelteKit web app, Wish/Bubble Tea TUI. Source of truth for product behavior: `docs/SPEC.md`. If code deviates from the spec, update the spec in the same PR.

## Non-negotiable: no secrets, no personal data

This is an open-source repo. Nothing tied to the maintainer's deployment or identity may ever be committed, pushed, put in an issue, PR, commit message, test fixture or log.

- All deployment-specific values come from **environment variables**: domains, hostnames, IPs, ports, Tailscale, Oracle/cloud accounts, registry credentials, SSH host key path, tokens, API keys. Code reads them via `internal/config` only.
- `.env.example` holds placeholders only (`example.com`, `changeme`). Real `.env*` files are gitignored.
- Deploy secrets live only in GitHub Actions secrets/environments. Workflows reference `${{ secrets.X }}`, never values.
- No real names, emails, company names, employer hostnames or corporate-network details anywhere. Use `example.com`, `Alex`, `Sam`, `team-checkout` in fixtures and docs. Refer to corporate proxies generically.
- Commit identity is `reh2ur <7092173+reh2ur@users.noreply.github.com>`. Never use or write any other email. Never change git identity config.
- Never print secret values in output; if gitleaks flags something, stop, report, do not commit.
- Never commit `*.pem`, `*_key`, `ssh_host_*`, `.env*` (except `.env.example`), data dumps.
- If a secret or personal data was committed, stop and tell the user; history rewrite and rotation are their call.

## Privacy rules from the spec

- The server stores nothing: no database, no disk writes of rooms or votes, no analytics, no third-party runtime assets.
- Only exception: the SSH host key file, read from `POKER_SSH_HOST_KEY_PATH`. In `POKER_ENV=production` the server must refuse to start without it. In dev an ephemeral key is generated with a warning.
- Logs never contain room names, display names or IP addresses.

## Repo layout

```
cmd/poker/          main: wires config, hub, web, ssh
internal/protocol/  versioned message types shared by web and ssh
internal/room/      hub, room state, reveal/result logic (pure, no I/O)
internal/web/       HTTP + WebSocket + go:embed of web/build
internal/sshd/      Wish + Bubble Tea, a client of the protocol
internal/config/    env-var config only
web/                SvelteKit + Tailwind (static adapter), pnpm
deploy/             Dockerfile, Caddyfile, compose
docs/               SPEC.md, self-hosting, ADRs
```

Rules: `room` imports nothing from `web` or `sshd`. `sshd` talks to the hub only through `protocol`, so the SSH gateway can be split out later.

## Commands

Use the Makefile: `make dev`, `make test`, `make lint`, `make build`. Toolchain is pinned: Go in `go.mod`, Node in `.nvmrc`, pnpm via `packageManager` in `web/package.json`.

## Workflow

- One issue = one branch = one PR. Branch `feat/12-room-hub`. PR body contains `Closes #12`. Squash merge.
- Conventional Commits: `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`, `ci:`.
- Milestones and labels (`area:web|ssh|core|infra|docs`, `type:feature|bug|spike|chore`, `owner:human`) organize work. Issues link to spec sections rather than copying text.
- Agents may push feature branches, open PRs and merge their own PR **only when all required CI checks are green** and the PR touches none of the protected paths below. Never push to `main` directly. Never force-push.
- Protected paths, which only the maintainer merges: `.github/` (incl. CODEOWNERS), `deploy/`, `.gitleaks.toml`, `lefthook.yml`, `CLAUDE.md`.
- Issues labeled `owner:human` are skipped by agents (infra provisioning, accounts, repo settings, company-network checks). Agents may prepare checklists, wizard scripts or commands for them but never execute them. Agents never change repo visibility, settings, branch protection or secrets.

## Quality gates

- Tests are required with every behavior change. Go: table-driven, `go test -race`. `internal/room` and `internal/protocol` keep at least 80% coverage (TDD preferred). Web: Vitest for stores and key handling.
- Never delete, skip or loosen tests, linters, coverage thresholds or CI checks to make a build pass. Fix the code.
- Before opening a PR run `make lint test`. CI also runs gitleaks, `go vet`, `golangci-lint`, `pnpm check`, build.
- Keep changes scoped to the issue. No speculative abstractions, no extra features.

## Style

- Default to no comments; add one only for a non-obvious why.
- Protocol messages are versioned from day one; join is idempotent via participant token.
- Keyboard-first: web and terminal share one interaction model (type a value, Enter commits).
