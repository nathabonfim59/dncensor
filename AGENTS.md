# dncensor — AGENTS.md

## Build & verify

```bash
make build       # release binary (enforces root)
make build-dev   # dev binary (skips root check)
make vet         # go vet ./...
make tidy        # go mod tidy
make clean       # remove binaries
make             # build + vet
```

No tests, no CI.

## Root requirements

- Root enforcement is in `cmd/requireRoot()` via build tags:
  - `cmd/rootcheck_release.go` (`//go:build !dev`) — checks `os.Geteuid() != 0` and exits.
  - `cmd/rootcheck_dev.go` (`//go:build dev`) — no-op for UI testing.
- `dncensor current` and `dncensor list-providers` work without root.
- `Stack.RequiresRoot()` exists on the interface but is **never called** in production code.

## Architecture

- **Entrypoint**: `main.go` → `cmd.Execute()` → root command (TUI) or subcommand.
- **Stack detection** (`internal/stack/stack.go:Detect()`): hardcoded order — systemd-resolved → NetworkManager → plain /etc/resolv.conf. First match wins.
- **ISP is not a real DNS provider**: `NewISP()` has no IPs. `dns.Apply()` special-cases `ProviderISP` to call `restoreISP()` (backup restore). No backup → error.
- **No automated restore**: All three `Restore()` methods return an error pointing to the backup dir. `restoreISP()` tries to use the file at `BackupRecord.BackupPath` but falls through to the same failing method.

## CLI flags

```
set --provider, -p  (required)  isp|cloudflare|google
    --flavor, -f                standard|malware|adult
    --doh                       enable DNS-over-HTTPS
    --yes, -y                   skip confirmation

current --json, -j
list-providers --json, -j
```

## Known quirks

1. **DoH overwrites DNS IPs on systemd-resolved**: `SetDOH()` calls `resolvectl dns <iface> <url>`, replacing the IP-based nameservers set by `SetDNS()`. DoH + IP DNS cannot coexist.
1. **Config paths are hardcoded** (`/etc/dncensor/`, `/etc/dncensor/backup/`). No env vars or flags to override.
1. **No config file** (no TOML/YAML/JSON). Only state is backup snapshots in `/etc/dncensor/backup/`.
1. **Backup files are timestamped** (no manifest, no metadata beyond filename prefix for stack type detection).

## Stack-specific notes

- **systemd-resolved**: `SetDNS` applies to every interface found in `resolvectl status`. Falls back to `lo` if parsing fails.
- **NetworkManager**: Prefers ethernet over wifi for active connection selection.
- **resolvconf**: Uses atomic write (temp file + rename). Preserves `search`/`domain` lines from existing config.

## Using external libraries

- **Find existing deps**: check `go.mod` and `go.sum`. When unsure if a library is already used, grep imports in the relevant `internal/` package.
- **Add a new dep**: `go get <module>@<version>` then `go mod tidy`. Keep it minimal — this project uses bubbletea, lipgloss, and cobra; there are no yaml/json parsers or test frameworks.
- **Understand usage**: read the existing callers in the codebase (e.g., `internal/tui/` for bubbletea patterns, `cmd/` for cobra) and run `go doc ./internal/... .Symbol` to inspect any exported API before writing new code.
