# ARCHITECTURE.md — Component Map

Living document. Update when package boundaries or responsibilities change.

**Deep dive:** [TECHNICAL.md](TECHNICAL.md)  
**Design invariants:** [.hermes/ARCHITECTURE_PRINCIPLES.md](.hermes/ARCHITECTURE_PRINCIPLES.md)

---

## System overview

```text
┌─────────────────────────────────────────────────────────────┐
│                        cmd/pork                             │
│  pork │ autopilot │ doctor │ proxy │ --version              │
└────────────┬────────────────────────────────────────────────┘
             │
    ┌────────┼────────┬──────────┬──────────┐
    ▼        ▼        ▼          ▼          ▼
 config   engine   aggregator  health     tui
    │        │        │          │          │
    │        │        └── provider         │
    │        │              rank           │
    │        ├── proxy                       │
    │        ├── state                       │
    │        └── isos (via TUI shelf)        │
    └────────────────────────────────────────┘
```

---

## Component responsibilities

| Component | Package | Owns | Does not own |
|-----------|---------|------|--------------|
| **CLI** | `cmd/pork` | Argument parsing, subcommand dispatch, wiring deps, background health warmup | UI rendering, provider HTML parsing |
| **TUI** | `internal/tui` | Screens (search, ISOs, results, preview, downloads, health), keys, styles, proxy badge | Torrent protocol, config file I/O |
| **Config** | `internal/config` | `~/.pork` paths, YAML load/save, defaults, proxy atomic updates | Search or download logic |
| **Aggregator** | `internal/aggregator` | Concurrent provider fan-out, retries, status events, content filter choke point | Per-site parsers |
| **Provider** | `internal/provider` | Index adapters, magnet resolution, response caps, NSFW tagging | Ranking, UI |
| **Rank** | `internal/rank` | Score formula, release tag parsing, noise detection | Fetching results |
| **Engine** | `internal/engine` | anacrolix client, add/pause/seed, direct HTTPS ISO downloads, snapshots, strict proxy transport policy | Search queries |
| **Proxy** | `internal/proxy` | SOCKS5 dialer, HTTP client factory, Tor egress check | Config persistence |
| **State** | `internal/state` | `state.json` entries, atomic save, upsert | Live torrent stats |
| **Health** | `internal/health` | Doctor checks, provider probes, swarm snapshots, trend store | TUI layout |
| **ISOs** | `internal/isos` | Distro catalog, live resolution, direct checksum verify | BitTorrent peer wire |
| **Autopilot** | `internal/autopilot` | Intent parse, batch select, headless progress (WIP) | Interactive search UX |

---

## Entry points

| Invocation | Primary packages |
|------------|------------------|
| `pork` | config → engine → state → aggregator → tui |
| `pork autopilot` | config → engine → state → aggregator → autopilot → (optional) tui |
| `pork doctor` | config (read-only) → health |
| `pork proxy` | config (proxy field only) |

---

## External dependencies

| Dependency | Used by | Role |
|------------|---------|------|
| `anacrolix/torrent` | engine | BitTorrent client |
| `charmbracelet/bubbletea` | tui | Terminal UI framework |
| `PuerkitoBio/goquery` | provider | HTML index scraping |
| `golang.org/x/net/proxy` | proxy | SOCKS5 dialer |

---

## Persistence map

| File | Writer | Reader |
|------|--------|--------|
| `~/.pork/config.yaml` | config.Load (first run), config.UpdateProxy, chmod hardening | config, doctor |
| `~/.pork/state.json` | state.Save (TUI exit, throttled ticks) | state.Load, engine resume |
| `~/.pork/health.json` | health.Store | health screen, doctor `--record` |
| Piece completion DB | engine (bolt under `~/.pork`) | engine resume |

---

## Screen flow (TUI)

```text
search ──tab──► isos ──tab──► results ──tab──► downloads
  │               │              │                 │
  │               │              ├── preview       └── health (H)
  └───────────────┴────────────── esc/back ────────────────┘
```

---

## When to update this file

- New `internal/` package or renamed boundary
- New CLI subcommand
- New persistence file or config section
- Provider added to default catalog
- TUI screen added or removed
