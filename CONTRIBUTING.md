# CONTRIBUTING.md — Pork

Thank you for contributing. This guide covers engineering standards, workflow, and quality gates for anyone working on pork.

For system design context, read [ARCHITECTURE.md](ARCHITECTURE.md) first.  
For deep technical reference, see [TECHNICAL.md](TECHNICAL.md).

---

## Engineering standards

All work on this repository is governed by **[`.hermes/`](.hermes/README.md)** — persistent engineering and documentation policies that travel with the repo.

Read before coding:

1. [.hermes/ENGINEERING_CONSTITUTION.md](.hermes/ENGINEERING_CONSTITUTION.md) — non-negotiable process rules
2. [.hermes/PROJECT_STANDARDS.md](.hermes/PROJECT_STANDARDS.md) — quality gates and safety rules
3. [.hermes/ARCHITECTURE_PRINCIPLES.md](.hermes/ARCHITECTURE_PRINCIPLES.md) — design invariants
4. [.hermes/DOCUMENTATION_STANDARD.md](.hermes/DOCUMENTATION_STANDARD.md) — doc structure

Before writing code:

1. Understand the real problem — not just the symptom
2. Design before implementing — boundaries, data flow, failure modes
3. Challenge your own design from multiple engineering perspectives
4. Implement progressively: correct first, optimize later
5. Self-review before submitting: races, leaks, API consistency, maintenance cost

Agents should start with [PROJECT_MANIFEST.md](PROJECT_MANIFEST.md).

---

## Who can contribute

- Bug fixes and documentation improvements are always welcome
- Feature work should start with an issue or design discussion for significant changes
- All contributors must follow the quality gates below

---

## Development setup

### Prerequisites

- Go 1.26 or newer
- macOS or Linux (primary targets)
- Optional: Nix for reproducible dev shell

### Clone and build

```bash
git clone https://github.com/maxiguillermo1/pork.git
cd pork
go build -o pork ./cmd/pork
./pork --version
```

### Nix dev shell

```bash
nix develop
go test ./...
```

### Run locally

```bash
go run ./cmd/pork
```

Config is created on first run at `~/.pork/config.yaml`. Downloads default to `~/Downloads/pork`.

---

## Workflow

1. **Understand** — Read ARCHITECTURE.md and the relevant `internal/` package
2. **Design** — For non-trivial changes, note tradeoffs in the PR or an ADR in TECHNICAL.md
3. **Implement** — Minimal focused diff; match existing conventions
4. **Test** — `go test -race`; add fixtures for provider HTML/JSON changes
5. **Verify** — Run all quality gates
6. **Document** — Update README, TECHNICAL, or ARCHITECTURE if behavior changes
7. **Changelog** — Add entry to CHANGELOG.md for user-visible changes

---

## Quality gates

Required before merge:

```bash
go vet ./...
go test -race -timeout 5m ./...
go build ./...
```

Run `go fmt` on files you edit.

---

## Pull request checklist

- [ ] [.hermes/ARCHITECTURE_PRINCIPLES.md](.hermes/ARCHITECTURE_PRINCIPLES.md) respected (especially proxy and ISO rules)
- [ ] Tests pass with race detector
- [ ] `go vet ./...` clean
- [ ] Documentation updated per [.hermes/PROJECT_STANDARDS.md](.hermes/PROJECT_STANDARDS.md) change-type table
- [ ] No secrets, credentials, or personal paths in the diff
- [ ] Provider changes include or update `testdata/` fixtures where practical
- [ ] CHANGELOG.md updated for user-visible changes

---

## Code conventions

- Keep packages under `internal/` — pork is not a public library
- Providers implement `provider.Provider`; do not parse HTML in the TUI
- Use `context.Context` for cancellable network work
- User-facing errors: calm, actionable, no raw stack traces
- Proxy credentials never logged or echoed in errors

---

## Packaging releases

When cutting a release, bump version consistently in:

- `packaging/homebrew/pork.rb`
- `packaging/nix/package.nix` and `flake.nix`
- `packaging/aur/PKGBUILD` and `.SRCINFO`
- `CHANGELOG.md`

Build with:

```bash
go build -ldflags "-s -w -X main.version=vX.Y.Z" -o pork ./cmd/pork
```

---

## Questions

Open a GitHub issue for bugs, provider outages, or design discussions.
