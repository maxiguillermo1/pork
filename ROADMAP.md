# ROADMAP.md

Planned and in-progress work for pork. Items reflect actual codebase status — nothing here is promised on a fixed date.

---

## In progress

### Autopilot

**Status:** WIP — functional but not finished

- [x] Intent parsing (`ParseIntent`) for queries, seasons, resolution hints
- [x] Batch gather, rank, and select (`Select`) with min-seeder floor
- [x] `--dry-run`, `-n`, `--headless` modes
- [ ] Smarter TV season / pack detection
- [ ] Interactive confirmation step before large queues
- [ ] Better error reporting when all providers fail
- [ ] Configurable default queries for ISO autopilot

---

## Near term

### Provider maintenance

Index sites change without notice. Ongoing work:

- [ ] Mirror rotation and fallback for 1337x and EZTV HTML scrapers
- [ ] Knaben/API health alerts in doctor output
- [ ] Document provider enable/disable recipes in README
- [ ] Expand `testdata/` coverage for regression tests when HTML changes
- [ ] Rate-limit awareness — backoff when providers return blocks

### CI improvements

- [ ] Run tests on macOS runner (darwin-specific path and proxy credential behavior)
- [ ] Cache Go module downloads across workflow runs (partially via `setup-go`)
- [ ] Lint job (`staticcheck` or `golangci-lint`) — evaluate signal vs noise
- [ ] Release workflow — tag, build artifacts, update packaging SHAs automatically
- [ ] SBOM or dependency review on release tags

---

## Future ideas

These are **not committed** — evaluate with ADR before building:

- Health history pruning / rotation policy
- Export/import config profiles
- Optional webhooks when downloads complete (headless/server use case)
- Plugin interface for providers (today: compile-time adapters only)
- Windows terminal polish and installer

---

## How to influence the roadmap

Open a GitHub issue with:

1. Problem you are solving
2. Whether you can contribute a PR
3. Whether it touches proxy, providers, or ISO policy (needs extra review)

See [CONTRIBUTING.md](CONTRIBUTING.md) for the PR process.
