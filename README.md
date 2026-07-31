# pork

```
  ^~~^
 ( ^.^ )  a cozy terminal for torrent search and official Linux ISOs
  > ᴥ <
```

**Pork** is a terminal app that searches public torrent indexes, ranks what it finds, and downloads files through BitTorrent — all from one calm screen. It also ships a curated shelf of **official Linux ISOs** so you never have to hunt sketchy mirror pages.

![pork home screen](docs/home.png)

---

## What is this?

Pork is a **BitTorrent client** with a built-in **search front end**. You type what you want, compare seeders and quality signals across several index sites at once, preview a torrent's files, and queue downloads without leaving the terminal.

Think of it as a small control room for torrents: search on the left, downloads on the right, and a health compass when you want to know whether your sources are still alive.

---

## Why was it built?

Most torrent tools assume you already have a magnet link or a `.torrent` file. Finding a good release across multiple sites — and comparing swarms — means juggling browser tabs, copy-pasting magnets, and hoping the indexer you bookmarked last month still works.

Pork puts **search**, **ranking**, **download management**, and **official ISO fetching** in one place, with optional **strict proxy** support for users who route traffic through Tor or SOCKS5.

---

## The problem

- Index sites differ in quality, speed, and availability; no single source is reliable.
- Duplicate releases clutter results; judging seeders, resolution, and source tags by eye is tedious.
- Linux ISOs should come from official publishers, but finding the current torrent or checksum on each project's site is repetitive.
- Privacy-conscious users need confidence that a proxy route is actually used — not silently bypassed when a peer protocol is inconvenient.

---

## The solution

Pork aggregates searches across enabled providers, scores and deduplicates results, and hands the winner to a BitTorrent engine wrapped in a friendly TUI. A separate **ISO shelf** resolves the latest official images from Ubuntu, Debian, Fedora, Arch, NixOS, Proxmox, and dozens more — live at selection time.

`pork doctor` checks your disk, config, providers, and (optionally) proxy egress. Health history stays local in `~/.pork/health.json`.

---

## How it works

1. **Launch** — `pork` loads `~/.pork/config.yaml`, resumes prior downloads from `state.json`, and starts the torrent engine.
2. **Search** — Your query fans out to every enabled provider in parallel (Knaben, YTS, Nyaa, plus optional RSS/Torznab feeds).
3. **Rank** — Results are scored by seeders, quality tags (1080p, WEB-DL, etc.), trust flags, and noise penalties.
4. **Preview & queue** — Pick a row, inspect files, then download. Magnets without a direct link are resolved lazily.
5. **Download** — BitTorrent pulls pieces; direct ISO downloads verify SHA-256 when the publisher publishes a checksum.
6. **Persist** — Progress and resume metadata save to `state.json` so restarts continue where you left off.

Press **H** in the TUI for the health compass. Press **Tab** to cycle screens.

---

## Technologies used

| Technology | What it is |
|------------|------------|
| **Go** | A compiled programming language known for simple deployment and strong concurrency. Pork is written in Go 1.26. |
| **Bubble Tea** | A Go framework for building terminal user interfaces (TUIs) — think "an app that runs inside your terminal window." |
| **Lipgloss** | A Go library for colors, borders, and layout in the terminal. Pork uses it for the cozy visual style. |
| **anacrolix/torrent** | A Go BitTorrent library. It handles peer connections, piece hashing, and resume data. |
| **YAML** | A human-readable config format. Pork stores settings in `~/.pork/config.yaml`. |
| **SOCKS5** | A proxy protocol. When enabled, pork routes TCP traffic through your proxy (for example local Tor on port 9050). |

---

## Software requirements

- **macOS or Linux** (primary targets; Windows may build but is less tested)
- **Go 1.26+** (only if building from source)
- A terminal with true-color support recommended

---

## Installation

**macOS (Homebrew)**

```sh
brew tap maxiguillermo1/tap && brew install pork
```

**Arch Linux (AUR)**

```sh
yay -S pork
```

**Nix**

```sh
nix run github:maxiguillermo1/pork
```

**Go install**

```sh
go install github.com/maxiguillermo1/pork/cmd/pork@latest
```

On macOS, if Go picks an unexpected compiler, use Apple Clang:

```sh
CC=/usr/bin/clang CXX=/usr/bin/clang++ go install github.com/maxiguillermo1/pork/cmd/pork@latest
```

---

## Usage

**Interactive (default)**

```sh
pork
```

Downloads land in `~/Downloads/pork` by default. Override:

```sh
pork -d /path/to/folder
```

**Diagnostics**

```sh
pork doctor
pork doctor --engine          # verify torrent listener
pork doctor --proxy-check     # verify proxy egress (Tor check service)
pork doctor --record          # save provider probe to health history
```

**Proxy (strict SOCKS5)**

```sh
pork proxy tor                # local Tor at 127.0.0.1:9050
pork proxy status
pork doctor --proxy-check
```

**Autopilot (work in progress)**

```sh
pork autopilot "ubuntu 24.04 desktop iso"
pork autopilot --dry-run "debian netinst amd64"
pork autopilot --headless -n 3 "fedora workstation x86_64"
```

**Evil mode** (session only — do not seed new downloads after they finish):

```sh
pork --evil
```

---

## Project structure

```text
cmd/pork/           Entry point and CLI subcommands
internal/           Application code (not importable by other modules)
  tui/              Terminal UI screens and keybindings
  engine/           BitTorrent and direct download engine
  provider/         Index site adapters
  aggregator/       Parallel search orchestration
  rank/             Result scoring
  config/           User configuration (~/.pork)
  state/            Download resume list
  health/           Doctor and health history
  isos/             Official Linux ISO catalog
  proxy/            SOCKS5 runtime
  autopilot/        Batch search and queue (WIP)
packaging/          Homebrew, Nix, and AUR recipes
docs/               Screenshots
.github/workflows/  CI (vet, build, race tests)
```

---

## Keyboard reference

| Screen | Keys |
|--------|------|
| **Home** | Type to search, `↑↓` pick destination, `Enter` go |
| **ISOs** | `↑↓` browse distros, `Enter` fetch latest official image |
| **Results** | `Enter` preview, `D` download now, `/` filter, `o` sort, `v` graph |
| **Downloads** | `p` pause/resume, `s` seed toggle, `v` verify, `m` move, `r` relink, `x` remove, `d` delete data |
| **Global** | `Tab` cycle screens, `Esc` back, `H` health, `Ctrl+C` quit |

---

## FAQ

**Where does pork store data?**  
Config and state live in `~/.pork/`. Downloads default to `~/Downloads/pork`.

**Which search providers are enabled by default?**  
Knaben, YTS, and Nyaa. Pirate Bay, EZTV, and 1337x are available but off by default — enable them in `config.yaml`.

**Does pork host or index content?**  
No. Pork is a client that queries third-party sites you configure. You are responsible for lawful use.

**What is strict proxy mode?**  
When a SOCKS5 proxy is set, pork disables protocols that could bypass it (DHT, uTP, UDP trackers, inbound peers). You may see fewer peers, but pork will not silently connect direct.

**Can I add my own indexer?**  
Yes — add an RSS or Torznab `search_url` under `providers` in `config.yaml`.

**Is autopilot finished?**  
Not yet. It can search, rank, and queue batches, but heuristics and UX are still evolving. See [ROADMAP.md](ROADMAP.md).

---

## Future roadmap

See [ROADMAP.md](ROADMAP.md) for planned work: autopilot polish, provider maintenance, and CI improvements.

---

## License

MIT License — Copyright (c) 2026 Maxi Guillermo. See [LICENSE](LICENSE).

Use pork only for content you are allowed to download and share. See [SECURITY.md](SECURITY.md) for proxy and credential guidance.
