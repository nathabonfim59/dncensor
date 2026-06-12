# dncensor

A TUI/CLI tool to detect your DNS stack and swap between DNS providers.

```
╭───────────────────────────────────────────────────────────────────╮
│                                                                   │
│  🛡 dncensor — DNS Provider Switcher                              │
│  Detected stack: systemd-resolved                                 │
│                                                                   │
│  Current DNS: Unknown                                             │
│                                                                   │
│  Select DNS Provider:                                             │
│   ▸ ○ ISP (restore original from DHCP)                            │
│     ○ CloudFlare                                                  │
│     ○ Google                                                      │
│                                                                   │
│     Apply Configuration                                           │
│                                                                   │
│  ☐ Use DNS-over-HTTPS (DoH)                                       │
│                                                                   │
│                                                                   │
│                                                                   │
│  ↑/↓ navigate · enter select · tab toggle DoH · a apply · q quit  │
│                                                                   │
╰───────────────────────────────────────────────────────────────────╯
```

## Why

I was having reliability issues with my ISP's DNS, and it was blocking
some sites. I use public DNS day-to-day but occasionally need to flip
back to my ISP's resolver for internal services. This makes the toggle
fast and safe.

## Features

- **Auto-detection:** figures out what DNS stack you are running
  (systemd-resolved, NetworkManager, resolvconf) and works with it.
- **Provider switching:** flip between ISP, Cloudflare, or Google DNS
  in one command.
- **Provider flavors:** Cloudflare offers standard, malware-blocking,
  and adult-filtering variants.
- **DNS-over-HTTPS:** toggle DoH on supported providers.
- **JSON output:** `current --json` and `list-providers --json` for
  programmatic consumption.
- **Named backups:** snapshot your current DNS config, then restore
  it later by name. Each backup is identified by its SHA256 hash.
- **TUI mode:** run `dncensor` with no arguments for an interactive
  terminal UI.

See `dncensor --help` and `dncensor <command> --help` for full usage.

## Build

```bash
make          # build + vet
make build    # release binary
make build-dev # dev binary (no root check)
```

Root is required at runtime for most commands (the dev build skips this
check for TUI testing). The release binary enforces it.
