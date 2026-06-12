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
- `dncensor backup list` requires root (reads `/etc/dncensor/snapshots/`).
- `Stack.RequiresRoot()` exists on the interface but is **never called** in production code.

## Architecture

- **Entrypoint**: `main.go` → `cmd.Execute()` → root command (TUI) or subcommand.
- **Stack detection** (`internal/stack/stack.go:Detect()`): hardcoded order — systemd-resolved → NetworkManager → plain /etc/resolv.conf. First match wins.
- **ISP detection** (`internal/dhcp/detect.go:DetectOriginalDNS()`): finds the default route interface, then looks for DNS servers in DHCP lease files (dhcpcd → dhclient → systemd-networkd). Used by `NewISP().Resolve()` to dynamically determine the ISP's original DNS.
- **User-driven backups** (`internal/backup/backup.go`): named snapshots with SHA256 hashes. No automatic backup — user explicitly creates with `dncensor backup create -n "name"`. Stored as JSON files in `/etc/dncensor/snapshots/<hash>.json`.
- **Stack interface** (`internal/stack/stack.go`): `CaptureDNS() ([]byte, error)` returns raw DNS config; `ApplyDNS(content []byte) error` applies it. No more `Backup`/`Restore` with hardcoded file paths.

## CLI flags

```
set --provider, -p  (required)  isp|cloudflare|google
    --flavor, -f                standard|malware|adult
    --doh                       enable DNS-over-HTTPS
    --yes, -y                   skip confirmation

current --json, -j
list-providers --json, -j

backup create --name, -n        backup name (prompted if empty)
backup list
backup restore <hash-or-name>   restore from backup
backup delete <hash-or-name>    delete a backup
```

## Known quirks

1. **DoH overwrites DNS IPs on systemd-resolved**: `SetDOH()` calls `resolvectl dns <iface> <url>`, replacing the IP-based nameservers set by `SetDNS()`. DoH + IP DNS cannot coexist.
1. **Config paths are hardcoded** (`/etc/dncensor/`, `/etc/dncensor/snapshots/`). No env vars or flags to override.
1. **No config file** (no TOML/YAML/JSON). Only state is user-created snapshots in `/etc/dncensor/snapshots/`.
1. **DHCP detection is read-only**: reads lease files from `/var/lib/dhcpcd/`, `/var/lib/dhcp/`, or `/run/systemd/netif/leases/`. The first client with a valid lease wins.

## Stack-specific notes

- **systemd-resolved**: `SetDNS` applies to every interface found in `resolvectl status`. Falls back to `lo` if parsing fails. `CaptureDNS` runs `resolvectl dns`, `ApplyDNS` parses that format and re-applies per-interface.
- **NetworkManager**: Prefers ethernet over wifi for active connection selection. `CaptureDNS` runs `nmcli -s -f ipv4.dns,ipv4.dns-search con show`, `ApplyDNS` re-applies via `nmcli con mod`.
- **resolvconf**: Uses atomic write (temp file + rename). `CaptureDNS` reads the resolved target file, `ApplyDNS` writes it back.

## Using external libraries

- **Find existing deps**: check `go.mod` and `go.sum`. When unsure if a library is already used, grep imports in the relevant `internal/` package.
- **Add a new dep**: `go get <module>@<version>` then `go mod tidy`. Keep it minimal — this project uses bubbletea, lipgloss, and cobra; there are no yaml/json parsers or test frameworks.
- **Understand usage**: read the existing callers in the codebase (e.g., `internal/tui/` for bubbletea patterns, `cmd/` for cobra) and run `go doc ./internal/... .Symbol` to inspect any exported API before writing new code.
