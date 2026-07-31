# PROJECT_MANIFEST.md — Agent Entry Point

**Purpose:** Machine-friendly project summary for AI coding agents. Humans should start with [README.md](README.md).

**Last verified:** 2026-07-30  
**Repository:** pork

---

## Identity

| Field | Value |
|-------|-------|
| **Name** | pork |
| **Domain** | Terminal BitTorrent search, download, and Linux ISO client |
| **Status** | Active — v0.2.0 packaged; autopilot WIP |
| **License** | MIT — Copyright (c) 2026 Maxi Guillermo |

---

## What this project does

Pork is a cozy terminal app for searching public torrent indexes, ranking results, downloading via BitTorrent (or verified HTTPS for select ISOs), and managing transfers in a Bubble Tea TUI. Optional strict SOCKS5 proxy mode routes all traffic through Tor or another proxy without silent fallback.

---

## Tech stack

| Layer | Technology |
|-------|------------|
| **Language** | Go 1.26 |
| **TUI** | Charm Bracelet Bubble Tea, Bubbles, Lipgloss |
| **BitTorrent** | anacrolix/torrent |
| **HTML/RSS parsing** | PuerkitoBio/goquery, custom RSS/Torznab |
| **Config** | YAML (`gopkg.in/yaml.v3`) in `~/.pork/` |
| **Build / deploy** | `go build`, Nix flake, Homebrew, AUR |
| **Test** | `go test -race`, table-driven tests, `testdata/` fixtures |

---

## Repository layout

```text
cmd/pork/              # main: TUI, autopilot, doctor, proxy CLI
internal/
  aggregator/          # concurrent multi-provider search
  autopilot/           # batch search + select + queue (WIP)
  config/              # ~/.pork config.yaml and paths
  engine/              # torrent + direct download engine
  health/              # doctor, probes, health.json history
  isos/                # official Linux ISO catalog
  provider/            # index adapters (Knaben, YTS, Nyaa, …)
  proxy/               # strict SOCKS5 runtime
  rank/                # result scoring and tag parsing
  state/               # state.json persistence
  tui/                 # Bubble Tea UI
packaging/             # homebrew, nix, aur
.github/workflows/     # CI: vet, build, race tests
docs/                  # screenshots for README
```

---

## Key commands

| Command | Purpose |
|---------|---------|
| `go mod download` | Fetch dependencies |
| `go run ./cmd/pork` | Run TUI locally |
| `go test -race -timeout 5m ./...` | Full test suite |
| `go vet ./...` | Static analysis |
| `go build -o pork ./cmd/pork` | Build binary |
| `go build -ldflags "-s -w -X main.version=dev" -o pork ./cmd/pork` | Release-style build |
| `nix develop` | Nix dev shell with Go |
| `pork doctor` | Read-only diagnostics |
| `pork proxy tor` | Enable local Tor SOCKS5 |
| `pork autopilot "query"` | Batch search and queue (WIP) |

---

## Architecture (summary)

CLI (`cmd/pork`) loads config and state, starts the torrent engine, builds the provider aggregator, and launches the Bubble Tea TUI (or subcommands). Providers fetch search results concurrently; rank scores them; the engine handles magnets and direct ISO downloads; state persists resume data; health records provider and swarm trends.

**Canonical detail:** [ARCHITECTURE.md](ARCHITECTURE.md) · [TECHNICAL.md](TECHNICAL.md)

---

## Data & configuration

| Item | Location |
|------|----------|
| **Config** | `~/.pork/config.yaml` |
| **Download state** | `~/.pork/state.json` |
| **Health history** | `~/.pork/health.json` |
| **Piece completion DB** | `~/.pork/` (bolt, managed by engine) |
| **Downloads default** | `~/Downloads/pork` (override: `-d`, `download_dir`) |
| **Secrets** | Proxy credentials in `config.yaml` only; never commit; mode `0600` required |

---

## Engineering standards

Read before editing:

1. [.hermes/README.md](.hermes/README.md) — policy index
2. [.hermes/PROJECT_STANDARDS.md](.hermes/PROJECT_STANDARDS.md) — quality gates
3. [.hermes/ARCHITECTURE_PRINCIPLES.md](.hermes/ARCHITECTURE_PRINCIPLES.md) — design invariants
4. [CONTRIBUTING.md](CONTRIBUTING.md) — workflow

---

## Testing workflow

```bash
go vet ./...
go test -race -timeout 5m ./...
go build ./...
```

Provider tests use files under `internal/provider/testdata/`. Engine integration tests avoid requiring live swarms where possible.

---

## Deployment

| Channel | Path |
|---------|------|
| Homebrew | `packaging/homebrew/pork.rb` |
| AUR | `packaging/aur/PKGBUILD` |
| Nix | `flake.nix`, `packaging/nix/package.nix` |
| Go install | `go install github.com/maxiguillermo1/pork/cmd/pork@latest` |

Version is injected via `-ldflags "-X main.version=vX.Y.Z"`.

---

## Constraints (do not violate)

- **Never** silently fall back from configured proxy to direct networking
- **Never** enable DHT/uTP/UDP/inbound peers in strict proxy mode
- **Never** store proxy credentials in shell history (use interactive `pork proxy set`)
- **Never** commit secrets, live provider credentials, or personal paths
- **Never** add ISO entries that are not official publisher infrastructure
- **Never** break the default `pork` TUI entry point or require env exports for normal use
- Provider adapters must respect `context` cancellation and response size caps
- `pork doctor` read-only path must not create or repair files

---

## Important files

| File | Why it matters |
|------|----------------|
| `cmd/pork/main.go` | CLI routing, provider wiring, health warmup |
| `internal/config/config.go` | Config schema, defaults, proxy updates |
| `internal/engine/engine.go` | Torrent client, strict proxy toggles |
| `internal/proxy/proxy.go` | SOCKS5 runtime and egress check |
| `internal/tui/app.go` | Screen flow and keybindings |
| `internal/provider/provider.go` | Provider contract |
| `internal/isos/isos.go` | ISO catalog and resolution |
| `internal/health/doctor.go` | Diagnostic report assembly |
| `flake.nix` | Nix package version source |

---

## Documentation map

| Audience | Start here |
|----------|------------|
| **Beginner** | [README.md](README.md) |
| **Contributor** | [CONTRIBUTING.md](CONTRIBUTING.md) |
| **Engineer** | [TECHNICAL.md](TECHNICAL.md) |
| **Agent** | This file |
