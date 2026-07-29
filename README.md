# pork

```
  ^~~^
 ( ^.^ )  a cozy terminal for torrent search + one-key linux isos
  > ᴥ <
```

You name it, the pig fetches it. Search public torrent indexes, compare the
actual swarm, and keep downloads in one calm little terminal app.

![pork home screen](docs/home.png)

## What it does

- Search Knaben, YTS, Nyaa, plus your own RSS or Torznab feeds.
- Group duplicate releases, rank the useful ones, and surface seeders, size,
  source, and noisy results before you download.
- Preview a magnet, queue it with one key, then pause, verify, seed, move, or
  relink it from the downloads screen.
- Browse and grab current official Linux ISOs from Ubuntu, Debian, Fedora,
  Arch, NixOS, Proxmox, and more. No sketchy mirror hunting.
- Run `pork doctor` when a provider or local setup feels off, and use the swarm
  compass to see whether your sources and downloads are healthy over time.

![pork search results](docs/results.png)

## Install

```sh
brew tap maxiguillermo1/tap && brew install pork             # macOS (Homebrew)
yay -S pork                                                  # Arch (AUR)
nix run github:maxiguillermo1/pork                           # Nix
go install github.com/maxiguillermo1/pork/cmd/pork@latest      # Go 1.26+
```

If Go on macOS is picking Nix's compiler, use Apple Clang for the install:

```sh
CC=/usr/bin/clang CXX=/usr/bin/clang++ go install github.com/maxiguillermo1/pork/cmd/pork@latest
```

Config lives in `~/.pork/`; downloads land in `~/Downloads/pork` (change with
`pork -d DIR`).

## SOCKS5 proxy

For the usual local Tor setup, one command is enough:

```sh
pork proxy tor
pork doctor --proxy-check
```

`pork proxy status` shows the redacted endpoint and strict-mode limits without
making a network request. To use another unauthenticated SOCKS5 proxy:

```sh
pork proxy set socks5://127.0.0.1:1080
```

For an authenticated endpoint, run `pork proxy set` without an argument and
enter the URL through hidden terminal input, so it never appears in shell
history or a process list. You can also configure pork by hand:

```yaml
proxy:
  socks5: "socks5://user:password@127.0.0.1:9050"
```

`socks5://` and `socks5h://` both resolve destination names at the proxy.
Username and password are optional and must be URL-encoded. For Tor, the usual
local endpoint is `socks5://127.0.0.1:9050`.

Credentials require `~/.pork/config.yaml` to be a regular file with mode
`0600`. A normal pork launch tightens an insecure regular config once; `pork
doctor` stays read-only and reports the problem instead.

Proxy mode is strict: searches, ISO requests, HTTP trackers, and outgoing TCP
peers use the proxy. pork disables DHT, uTP, UDP trackers, inbound peers, port
forwarding, and WebTorrent rather than risk a direct connection. Downloads may
find fewer peers as a result. The TUI keeps a visible `SOCKS strict` badge; it
becomes `((o)) Tor strict` only after a live verification confirms a Tor exit. A
check-service outage is shown as unavailable, not as a leak, and pork never
falls back to a direct connection.

## Keep an eye on it

`pork doctor` is read-only by default and checks config, disk, state, and
provider reachability. Add `--engine` for an opt-in listener check, `--record`
to save provider results in `~/.pork/health.json`, or `--proxy-check` to ask the
Tor Project's check service what egress IP it sees through your configured
proxy. That proves this explicit HTTP check used the route; it is not a promise
of anonymity or a claim about your normal connection. Health history contains
local provider timings plus torrent names and swarm counts; it is never
uploaded.

Automatic checks are off by default. Enable a local daily check with:

```yaml
health:
  enabled: true
  interval_hours: 24
```

The check sends the generic `1080p` canary query to enabled providers. Press
`H` in pork to view saved source and swarm history, or `r` there to record a
manual check.

## Keys

- **home** type to search, `↑↓` pick a destination, `enter` go
- **isos** `↑↓` browse, `enter` grab the latest official image
- **results** `enter` preview/get, `D` grab now, `/` filter, `o` sort, `v` graph
- **downloads** `p` pause/resume, `s` seed, `v` verify, `m` move, `r` relink, `x` remove, `d` delete data
- `tab` cycle, `esc` back, `^c` quit

## Autopilot (WIP)

Describe what you want and let the pig fetch it:

```sh
pork autopilot "all breaking bad seasons 1080p"      # also: --dry-run, -n N, --headless
```

## Legal

pork is a BitTorrent client and search tool. It does not host files, operate
trackers, or control the third-party providers it can search.

Use it only for content you are allowed to download and share, such as official
Linux ISOs, public-domain media, open-source software, and your own files. You
are responsible for following local law and the terms of any provider you
enable.

Provider availability and results can change without notice. pork does not
endorse or guarantee third-party content.

A proxy routes pork's traffic, but it is not a promise of anonymity or legal
protection.

MIT — Copyright (c) 2026 Maxi Guillermo. See [LICENSE](LICENSE).
