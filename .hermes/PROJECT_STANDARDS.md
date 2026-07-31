# Project Standards — Pork

Version: 1.0  
Repository: pork

---

## Governance hierarchy

1. **[.hermes/ENGINEERING_CONSTITUTION.md](ENGINEERING_CONSTITUTION.md)** — non-negotiable engineering rules
2. **[DOCUMENTATION_STANDARD.md](DOCUMENTATION_STANDARD.md)** — documentation requirements
3. **[ARCHITECTURE_PRINCIPLES.md](ARCHITECTURE_PRINCIPLES.md)** — design invariants
4. **[CONTRIBUTING.md](../CONTRIBUTING.md)** — contributor workflow and PR checklist

When policies conflict on repo-specific concerns, ARCHITECTURE_PRINCIPLES.md wins for runtime behavior; ENGINEERING_CONSTITUTION.md wins for process.

---

## Immutable user workflow

```text
1. Install pork (brew, AUR, Nix, or go install).
2. Run: pork
3. Search, browse ISOs, or queue downloads from the TUI.
```

Never require manual environment exports for normal use. The TUI must remain the default entry point.

---

## Critical safety rules

| Rule | Rationale |
|------|-----------|
| Proxy mode must **never** silently fall back to direct networking | User trusts the configured privacy route |
| Strict SOCKS5 disables DHT, uTP, UDP trackers, inbound peers, and WebTorrent | Prevent accidental direct connections |
| Credentialed proxy URLs require `config.yaml` mode `0600` | Protect secrets on disk |
| Refuse symlinked or non-regular config for mutating proxy edits | Prevent path escape on atomic rename |
| ISO shelf serves **official** images from **official** infrastructure only | Lawful-content charter |
| Do not commit secrets, personal paths, or provider credentials | Security and privacy |
| Provider scrapers are fragile — test with fixtures, not live sites in CI | Reliable CI and respectful crawling |

---

## Quality gates (required before merge)

```bash
go vet ./...
go test -race -timeout 5m ./...
go build ./...
```

Optional local check (matches release packaging):

```bash
go build -ldflags "-s -w -X main.version=dev" -o pork ./cmd/pork
```

---

## Code style

- Match existing Go package boundaries (`cmd/`, `internal/`)
- Read [ARCHITECTURE.md](../ARCHITECTURE.md) and the relevant `internal/` package before editing
- Design before implementing (see engineering constitution lifecycle)
- Verify claims against code — do not trust README alone
- Minimal scope: smallest correct diff
- Run `go fmt` on touched files before commit

---

## Documentation requirements per change type

| Change type | Required doc updates |
|-------------|---------------------|
| New CLI command or flag | README.md, TECHNICAL.md, CHANGELOG.md, PROJECT_MANIFEST.md |
| Architecture change | ARCHITECTURE.md, TECHNICAL.md (+ ADR if significant) |
| New package under `internal/` | ARCHITECTURE.md, TECHNICAL.md |
| Provider added or removed | README.md, TECHNICAL.md, CHANGELOG.md |
| Proxy or security change | SECURITY.md, TECHNICAL.md, README.md |
| Config schema change | README.md, TECHNICAL.md, CHANGELOG.md |
| Breaking change | CHANGELOG.md, README.md |

---

## Pull request checklist

- [ ] Architecture principles and safety rules respected
- [ ] Tests pass (`go test -race -timeout 5m ./...`)
- [ ] Vet clean (`go vet ./...`)
- [ ] Build succeeds (`go build ./...`)
- [ ] Documentation updated per change type table above
- [ ] No secrets, credentials, or personal paths in committed files
- [ ] Provider changes include or update `testdata/` fixtures where practical

---

## Packaging

Release versions are set in:

- `packaging/homebrew/pork.rb`
- `packaging/nix/package.nix` and `flake.nix`
- `packaging/aur/PKGBUILD`

Bump all packaging files together when cutting a release. Linker flag: `-X main.version=vX.Y.Z`.
