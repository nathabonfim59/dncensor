# dncensor — AGENTS.md

## Build & verify

```bash
go build -o dncensor .
go vet ./...
go mod tidy            # after dep changes
```

No tests, no CI, no Makefile.

## Root requirements

- `dncensor` (TUI) and `dncensor set` check `os.Geteuid() != 0` at the cobra layer and exit.
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

1. **Flavor menu hardcoded to index 1** (`internal/tui/main.go:174`): `m.selectedIdx == 1`, not `p.SupportsFlavors()`. If provider order changes, CloudFlare flavor menu breaks.
2. **DoH overwrites DNS IPs on systemd-resolved**: `SetDOH()` calls `resolvectl dns <iface> <url>`, replacing the IP-based nameservers set by `SetDNS()`. DoH + IP DNS cannot coexist.
3. **Config paths are hardcoded** (`/etc/dncensor/`, `/etc/dncensor/backup/`). No env vars or flags to override.
4. **No config file** (no TOML/YAML/JSON). Only state is backup snapshots in `/etc/dncensor/backup/`.
5. **Backup files are timestamped** (no manifest, no metadata beyond filename prefix for stack type detection).

## Stack-specific notes

- **systemd-resolved**: `SetDNS` applies to every interface found in `resolvectl status`. Falls back to `lo` if parsing fails.
- **NetworkManager**: Prefers ethernet over wifi for active connection selection.
- **resolvconf**: Uses atomic write (temp file + rename). Preserves `search`/`domain` lines from existing config.

## Dependencies

- `bubbletea` v1.3 — TUI (older `tea.Model` interface, no `tea.Program` options)
- `lipgloss` v1.1 — terminal styling
- `cobra` v1.10 — CLI framework
- No yaml/json parsers beyond stdlib, no test deps.
