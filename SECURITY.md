# SECURITY.md

Security practices, lawful use, and proxy guidance for pork.

**Technical detail:** [TECHNICAL.md](TECHNICAL.md) §5 (Proxy model)  
**Architecture invariants:** [.hermes/ARCHITECTURE_PRINCIPLES.md](.hermes/ARCHITECTURE_PRINCIPLES.md)

---

## Reporting vulnerabilities

If you believe you have found a security issue in pork, please report it responsibly:

1. **Do not** open a public GitHub issue for undisclosed vulnerabilities
2. Contact the maintainer via GitHub private security advisory or the email listed on the maintainer profile
3. Include steps to reproduce, affected version, and impact assessment

---

## Lawful use

Pork is a **BitTorrent client and search tool**. It does not host files, operate trackers, or control third-party index sites.

**You are responsible for:**

- Complying with copyright and licensing laws in your jurisdiction
- Following the terms of service of any index or provider you enable
- Sharing only content you have the right to distribute

**Intended lawful uses include:**

- Downloading official Linux ISOs from publisher infrastructure (pork's ISO shelf is designed for this)
- Open-source software distributed via BitTorrent
- Public-domain media
- Your own files

Pork does not endorse or verify third-party search results. Enabling additional providers (TPB, 1337x, etc.) is opt-in.

---

## Proxy security model

### What strict proxy mode does

When `proxy.socks5` is set in `~/.pork/config.yaml`:

- Provider HTTP searches use the SOCKS5 route
- ISO resolution HTTP requests use the SOCKS5 route
- Outgoing BitTorrent TCP peer connections use the SOCKS5 route
- HTTP(S) trackers and webseeds use the SOCKS5 route
- DHT, uTP, UDP trackers, inbound peer listeners, port forwarding, and WebTorrent are **disabled**

### What it does not guarantee

- **Anonymity** — A proxy routes pork's traffic; it does not make you anonymous by itself. Tor usage has operational security requirements beyond software settings.
- **DNS privacy** — Use `socks5h://` (or `socks5://` with remote resolution) so hostnames resolve at the proxy. Pork accepts both schemes and sends destination hostnames to the proxy.
- **Non-pork traffic** — Only pork's process is affected. Other applications on your system are unchanged.
- **Legal protection** — Routing through Tor or another proxy does not exempt you from applicable law.

### No silent fallback

If proxy configuration is invalid, pork **refuses to start** on the normal load path rather than connecting direct. This is intentional.

If the Tor Project check service is unreachable, pork reports a **verification outage** — not a leak. Route failures are distinguished from check-service failures in code (`proxy.IsRouteFailure`).

---

## Credential handling

### Storage

Proxy credentials may appear in `proxy.socks5` as URL userinfo:

```yaml
proxy:
  socks5: "socks5://user:password@127.0.0.1:9050"
```

**Requirements:**

- `~/.pork/config.yaml` must be a regular file (not a symlink) for mutating edits
- On Unix, mode `0600` (owner read/write only) is enforced when credentials are present
- Pork auto-chmods on normal launch if needed; `pork doctor` reports insecure permissions read-only

### CLI entry

- **Unauthenticated proxy:** `pork proxy set socks5://127.0.0.1:1080`
- **Authenticated proxy:** run `pork proxy set` with no URL argument — enter the URL via hidden terminal input (not shell history, not argv)

Never pass credentialed URLs on the command line.

### Logging and diagnostics

- Proxy errors do not echo the configured URL or password
- `pork proxy status` prints a **redacted** endpoint (host:port only)
- Doctor output does not include proxy credentials

---

## Local data sensitivity

| File | Contents | Exposure risk |
|------|----------|---------------|
| `~/.pork/config.yaml` | Settings, optional proxy creds | Protect file permissions |
| `~/.pork/state.json` | Torrent names, magnets, paths | Reveals download history |
| `~/.pork/health.json` | Provider timings, swarm stats, torrent names | Local diagnostic only — never uploaded by pork |
| Download directory | Downloaded content | User-controlled |

Pork does not phone home. Health and diagnostic data stay on disk unless you share them.

---

## Network surface

| Mode | Inbound listener | UDP | Notes |
|------|------------------|-----|-------|
| Direct | Ephemeral TCP port | DHT/uTP enabled | Standard BitTorrent client surface |
| Strict proxy | None | Disabled | Reduced peer discovery |

Use a firewall as appropriate for your threat model. In strict proxy mode, pork does not open an inbound torrent port.

---

## Provider and scraping risks

- Providers fetch third-party websites; pork caps response size (32 MiB) to limit memory exhaustion
- Scraped HTML is parsed in-process; malformed responses should fail the provider attempt, not crash the app
- `ErrBlocked` indicates anti-bot protection — pork does not attempt credential stuffing or bypass

Do not point RSS/Torznab URLs at untrusted hosts you do not control.

---

## Dependency hygiene

- Go modules pinned in `go.sum`
- CI runs `go vet` and race-enabled tests on push/PR
- Review dependency updates for torrent and networking libraries (`anacrolix/torrent`, `golang.org/x/net`)

---

## Secure development practices

Contributors should:

- Never commit real proxy credentials or personal `config.yaml` snippets
- Add provider fixtures under `testdata/` instead of live credentials in tests
- Run `go test -race` before submitting concurrency changes
- Read [.hermes/PROJECT_STANDARDS.md](.hermes/PROJECT_STANDARDS.md) before proxy or provider changes

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full PR checklist.
