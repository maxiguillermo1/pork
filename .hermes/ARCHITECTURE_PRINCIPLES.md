# Architecture Principles — Pork

Version: 1.0  
Living document — update when core invariants change.  
Complements [ARCHITECTURE.md](../ARCHITECTURE.md) (component map) and [TECHNICAL.md](../TECHNICAL.md) (implementation detail).

---

## Core purpose

Pork is a **terminal BitTorrent client** with integrated **index search**, **result ranking**, and a curated **official Linux ISO shelf**. It wraps `anacrolix/torrent` behind a Bubble Tea TUI and optional strict SOCKS5 routing.

---

## Design invariants

These must never be violated without an explicit ADR and user-visible migration path.

### 1. No silent proxy fallback

When `proxy.socks5` is configured, every outbound TCP path — provider searches, ISO resolution, HTTP trackers, and peer connections — must use the proxy. Invalid proxy configuration **blocks startup** on the mutating load path. Pork never “helpfully” retries direct.

### 2. Strict proxy is strict

In proxy mode, pork disables DHT, uTP, UDP trackers, inbound listeners, port forwarding, and WebTorrent. Fewer peers is acceptable; accidental direct egress is not.

### 3. Provider fragility is expected

Third-party index sites change HTML, mirrors, and anti-bot rules without notice. Providers are isolated adapters behind a narrow `Provider` interface. The aggregator owns retries and timeouts; the TUI never parses HTML.

### 4. Lawful ISO charter

The ISO shelf (`internal/isos`) resolves **official** release images from **official** publisher infrastructure — torrents where available, verified HTTPS direct downloads where not. Pork does not mirror arbitrary third-party ISOs or endorse unverified images.

### 5. Fail soft on resume, fail loud on privacy

Corrupt `state.json` or `config.yaml` on the normal load path is backed up and reset so the app can start. Proxy misconfiguration on the normal load path **refuses** to start. `pork doctor` is read-only and reports without mutating.

### 6. One choke point for content policy

NSFW filtering and provider result shaping flow through `aggregator.WithFilter` so TUI, autopilot, and future callers share identical rules.

### 7. Measured, not estimated

Performance or health claims in documentation must reflect code behavior or recorded `health.json` snapshots — not guesses.

---

## Separation of concerns

| Layer | Responsibility |
|-------|----------------|
| `cmd/pork` | CLI dispatch: TUI, autopilot, doctor, proxy |
| `internal/tui` | Bubble Tea screens, keys, rendering |
| `internal/aggregator` | Concurrent multi-provider search |
| `internal/provider` | Per-site search and magnet resolution |
| `internal/rank` | Scoring, tag parsing, noise detection |
| `internal/engine` | Torrent client, direct HTTPS downloads, snapshots |
| `internal/proxy` | SOCKS5 runtime, egress verification |
| `internal/config` | `~/.pork` paths and YAML |
| `internal/state` | Persisted download list (`state.json`) |
| `internal/health` | Doctor, probes, trend storage |
| `internal/isos` | Curated distro catalog and live resolution |
| `internal/autopilot` | Intent parsing, batch select, queue |

Modules must not leak responsibilities across these boundaries.

---

## Data flow principles

```text
User types query in TUI
    → aggregator fans out to enabled providers (shared HTTP client / proxy)
    → results stream back; rank scores and deduplicates
    → user previews magnet → engine adds torrent
    → state.json persists resume metadata
    → engine snapshots drive downloads UI and health probes
```

ISO shelf flow:

```text
User picks distro
    → isos resolves latest official image (live scrape/API)
    → engine adds magnet or direct HTTPS download with checksum verify
```

---

## Error handling principle

- User-facing errors are calm and actionable (autopilot `ShortReason`, doctor glyphs)
- Background health checks log to stderr but never crash the TUI
- Provider failures are per-provider; one dead mirror must not block others
- Doctor exits non-zero on hard failures for cron/canary use

---

## Testing principle

- Unit tests with `testdata/` fixtures for HTML/JSON/RSS providers
- Race detector enabled in CI (`go test -race`)
- Integration tests for engine and config without live network where possible
- Live network tests gated behind build tags or skipped in CI

---

## Future scaling considerations

- Autopilot is WIP — selection heuristics will evolve; keep intent parsing separate from queue logic
- New providers → new adapter + fixture tests + config default entry (disabled if fragile)
- Tor verification badge → depends on third-party check service availability; distinguish route failure from service outage
- Remote headless autopilot may grow; preserve shared `aggregator` + `engine` core
