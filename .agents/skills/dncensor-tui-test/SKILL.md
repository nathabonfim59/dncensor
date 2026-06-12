---
name: dncensor-tui-test
description: Test the dncensor TUI using tmux. Never run the app from the agent — always use a tmux pane so you can see the UI and interact properly.
---

# dncensor TUI test via tmux

## When to use

Visual/interactive testing of dncensor's TUI. Never run the binary from the agent CLI — use tmux.

## Steps

1. **Check tmux** — if `$TMUX` unset, tell user and stop.

2. **Pick or create window**:
   ```bash
   WIN=dncensor-test
   tmux list-windows -F '#{window_name}' | grep -qx "$WIN" \
     || tmux new-window -n "$WIN" -c "$PWD"
   ```

3. **Build dev binary**:
   ```bash
   make build-dev
   ```

4. **Send keys** — send Enter separately from command (avoids buffering races):
   ```bash
   tmux send-keys -t "$WIN" './dncensor-dev' Enter
   ```

5. **Capture visible output**:
   ```bash
   tmux capture-pane -t "$WIN" -p
   ```

6. **Navigate TUI**:
   ```bash
   tmux send-keys -t "$WIN" Down     # move down
   tmux send-keys -t "$WIN" Up       # move up
   tmux send-keys -t "$WIN" Enter    # select
   tmux send-keys -t "$WIN" Tab      # toggle DoH
   tmux send-keys -t "$WIN" 'a'      # apply
   tmux send-keys -t "$WIN" 'q'      # quit
   ```

7. **Subcommands** (no TUI):
   ```bash
   tmux send-keys -t "$WIN" './dncensor-dev current --json' Enter
   tmux send-keys -t "$WIN" './dncensor-dev list-providers --json' Enter
   ```

## Notes

- `make build-dev` uses `//go:build dev` tag to skip root check → binary is `./dncensor-dev`.
- Send `C-c` first if the pane has leftover input before sending fresh commands.
- Always use `-p` on capture-pane (visible only), never `-S -` (scrollback).
- The TUI renders inside the pane — capture output includes padding/whitespace from tmux; that's normal.
