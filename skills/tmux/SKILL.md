---
name: tmux
description: Guide for using tmux to run and manage long-running processes programmatically. Use when AI agents need to start background processes, capture their output, and manage them via the tmux CLI. Ideal for dev servers, build processes, and other commands that need to run independently of the agent's lifecycle.
---

# TMUX for Programmatic Process Management

TMUX allows agents to run long-running commands in background sessions that persist beyond the agent's execution. Sessions can be queried, controlled, and managed via CLI commands.

## Core Commands

### Starting Background Processes

Create a detached session running a command:

```bash
tmux new-session -d -s <session-name> "<command>"
```

Examples:

```bash
tmux new-session -d -s devserver "npm run dev"
```

```bash
tmux new-session -d -s build "make build"
```

```bash
tmux new-session -d -s python "python3 server.py"
```

**Best practices for session names:**

- Use descriptive names: `devserver`, `test-runner`, `watch-build`
- Avoid spaces and special characters
- Consider prefixing with project name if managing multiple projects

### Sending Commands to Sessions

Send input to a running session:

```bash
tmux send-keys -t <session-name> "<command>" C-m
```

The `C-m` sends a carriage return (equivalent to pressing Enter).

Examples:

```bash
# Send a restart message
tmux send-keys -t devserver "echo 'restarting...'" C-m
```

```bash
# If 'rs' is a shell alias for restart
tmux send-keys -t devserver "rs" C-m
```

### Capturing Output

Capture the current pane content:

```bash
tmux capture-pane -t <session-name> -p
```

Options:

- `-p`: Print to stdout (instead of internal buffer)
- `-S -`: Capture entire scrollback history
- `-e`: Include escape sequences (for colors)
- `-N <lines>`: Capture specific number of lines from bottom

Examples:

```bash
# Capture visible content
tmux capture-pane -t devserver -p
```

```bash
# Capture full history
tmux capture-pane -t devserver -S - -p
```

```bash
# Capture last 50 lines
tmux capture-pane -t devserver -N 50 -p
```

### Listing Sessions

```bash
tmux list-sessions
# or
tmux ls
```

Output format: `<session-name>: <windows> (created <timestamp>) [attached]`

### Checking Session Status

Check if a specific session exists:

```bash
tmux has-session -t <session-name>
```

Returns exit code 0 if exists, 1 if not. Useful in scripts.

### Killing Sessions

Terminate a session:

```bash
tmux kill-session -t <session-name>
```

Kill all sessions:

```bash
tmux kill-server
```

**Caution:** `kill-server` terminates all tmux sessions system-wide.

## Advanced Usage

### Multiple Windows in a Session

Create a new window in an existing session:

```bash
tmux new-window -t <session-name> -n <window-name> "<command>"
```

List windows in a session:

```bash
tmux list-windows -t <session-name>
```

### Targeting Specific Panes

By default, commands target pane 0 of the current window. To target specific panes:

```bash
tmux send-keys -t <session-name>:<window>.<pane> "<command>" C-m
tmux capture-pane -t <session-name>:<window>.<pane> -p
```

Example:

```bash
# Captures pane 2 of window 1
tmux capture-pane -t devserver:1.2 -p
```

### Environment Variables

Set environment variables for a session:

```bash
tmux new-session -d -s mysession "ENV_VAR=value npm run dev"
```

## Common Patterns

### Pattern 1: Start, Monitor, Cleanup

```bash
# Start a dev server
tmux new-session -d -s devserver "npm run dev"

# Wait a bit, then check output
sleep 5
tmux capture-pane -t devserver -p

# Later, when done
tmux kill-session -t devserver
```

### Pattern 2: Conditional Start

```bash
# Only start if session doesn't exist
if ! tmux has-session -t devserver; then
    tmux new-session -d -s devserver "npm run dev"
fi
```

### Pattern 3: Run Command, Get Output

```bash
# Send command and immediately capture output
tmux send-keys -t devserver "npm test" C-m
sleep 10
tmux capture-pane -t devserver -N 100 -p | grep -A 20 "Tests:"
```

## Troubleshooting

### Session Already Exists

If creating a session that exists, tmux will error. Check first:

```bash
tmux has-session -t <name> && echo "exists" || echo "free"
```

### No Output in capture-pane

Ensure the session is actually producing output. Some processes buffer output unless they detect a TTY. Forcing PTY allocation is not directly available via tmux CLI.

### Session Not Responding

Check if the process hung inside the session. Use `tmux capture-pane` to see the last output. If necessary, use `kill-session` to terminate.

## Configuration

For programmatic use, consider creating `~/.tmux.conf`:

```
# Disable status bar for cleaner output capture
set -g status off

# Increase scrollback buffer for more history
set -g history-limit 10000
```

This is optional but helps when capturing large amounts of output.
