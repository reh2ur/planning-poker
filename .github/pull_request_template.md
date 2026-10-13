Closes #

## Summary

<!-- What changed and why. Link spec sections rather than copying text. -->

## Test notes

<!-- Tests added or changed, and the result of `make lint test`. -->

## Checklist

- [ ] Scoped to the linked issue, no extra features
- [ ] Tests added or updated for behavior changes; `make lint test` passes
- [ ] `docs/SPEC.md` updated if behavior deviates from it
- [ ] No secrets or personal data in code, fixtures, logs, comments, commit messages or this PR (see CLAUDE.md)
- [ ] No real names, emails, company names, hostnames or IP addresses; placeholders only (`example.com`, `Alex`, `team-checkout`)
- [ ] Deployment-specific values come from environment variables via `internal/config`
- [ ] Logs contain no room names, display names or IP addresses
- [ ] Touches a protected path (`.github/`, `deploy/`, `.gitleaks.toml`, `lefthook.yml`, `CLAUDE.md`): maintainer merges
