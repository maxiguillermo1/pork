# CHANGELOG.md

All notable changes to pork are documented here.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).  
Versioning follows packaging releases (`v0.1.0`, `v0.2.0`).

---

## [0.2.0] — 2026-07-28

### Added

- **Strict SOCKS5 proxy mode** — route provider searches, ISO requests, HTTP trackers, and TCP peers through a configured proxy (`pork proxy tor|set|off|status`)
- **Proxy egress verification** — `pork doctor --proxy-check` and in-TUI Tor badge via Tor Project check service
- **Credential-safe proxy setup** — interactive `pork proxy set` for authenticated URLs; config mode `0600` enforcement
- **Engine hardening for proxy** — disables DHT, uTP, UDP trackers, inbound listeners, and WebTorrent in strict mode
- **Health and doctor improvements** — read-only doctor path, optional `--engine` probe, `--record` for health history
- **Download progress UI polish** — cleaner progress display in downloads screen
- **Torrent file preview** — extension badges on file tree preview

### Changed

- Rebranded from **tork** to **pork** (binary, module path, packaging)
- README and install docs updated for pork naming
- Nix flake and packaging recipes aligned to v0.2.0

### Removed

- Experimental automatic MKV→MP4 remux (added in development, reverted before release)

### Security

- Invalid proxy configuration blocks normal startup (no silent direct fallback)
- Proxy URL validation never echoes credentials in error messages
- Symlinked config refused for mutating proxy edits

---

## [0.1.0] — 2026-07

### Added

- **Initial release** — terminal torrent search and download TUI
- **Multi-provider search** — Knaben, YTS, Nyaa (defaults); optional TPB, EZTV, 1337x, RSS/Torznab
- **Result ranking** — seeders, quality tags, trusted releases, noise filtering
- **Downloads dashboard** — pause, resume, seed toggle, verify, move, relink, remove
- **Official Linux ISO shelf** — live resolution for major distros (Ubuntu, Debian, Fedora, Arch, NixOS, Proxmox, and more)
- **Direct HTTPS ISO downloads** — verified checksums for distros without official torrents
- **Autopilot (preview)** — `pork autopilot` batch search and queue
- **Health compass** — local `health.json` history; `pork doctor` diagnostics
- **Packaging** — Homebrew, AUR, Nix flake
- **CI** — `go vet`, build, and race-enabled tests on push/PR

### Notes

- MIT License — Copyright (c) 2026 Maxi Guillermo

[0.2.0]: https://github.com/maxiguillermo1/pork/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/maxiguillermo1/pork/releases/tag/v0.1.0
