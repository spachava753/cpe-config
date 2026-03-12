# Local Setup

Keep local browser configuration in shell environment variables rather than hard-coding personal paths in the skill.

## Expected Environment Variables

Required:
- `CPE_BROWSER_CHROME_BIN`
- `CPE_BROWSER_AUTOMATION_USER_DATA_DIR`
- `CPE_BROWSER_REMOTE_DEBUG_HOST`
- `CPE_BROWSER_VISIBLE_DEBUG_PORT`
- `CPE_BROWSER_HEADLESS_DEBUG_PORT`
- `CPE_BROWSER_VISIBLE_DEBUG_URL`
- `CPE_BROWSER_HEADLESS_DEBUG_URL`

Recommended:
- `CPE_BROWSER_ARTIFACT_DIR`
- `CPE_BROWSER_VIEWPORT_WIDTH`
- `CPE_BROWSER_VIEWPORT_HEIGHT`
- `CPE_BROWSER_DEVICE_SCALE_FACTOR`

Fail early with a clear error if required variables are missing.

## Launch Strategy

The agent should launch Chrome directly when needed.

Before launching, check whether `CPE_BROWSER_VISIBLE_DEBUG_URL` or `CPE_BROWSER_HEADLESS_DEBUG_URL` already responds and reuse that process when appropriate.

Visible mode launch pattern:
```text
$CPE_BROWSER_CHROME_BIN \
  --user-data-dir="$CPE_BROWSER_AUTOMATION_USER_DATA_DIR" \
  --remote-debugging-address="$CPE_BROWSER_REMOTE_DEBUG_HOST" \
  --remote-debugging-port="$CPE_BROWSER_VISIBLE_DEBUG_PORT"
```

Headless mode launch pattern:
```text
$CPE_BROWSER_CHROME_BIN \
  --user-data-dir="$CPE_BROWSER_AUTOMATION_USER_DATA_DIR" \
  --remote-debugging-address="$CPE_BROWSER_REMOTE_DEBUG_HOST" \
  --remote-debugging-port="$CPE_BROWSER_HEADLESS_DEBUG_PORT" \
  --headless=new
```

Apply the configured viewport after connecting, using `CPE_BROWSER_VIEWPORT_WIDTH`, `CPE_BROWSER_VIEWPORT_HEIGHT`, and `CPE_BROWSER_DEVICE_SCALE_FACTOR`.

## Tiny Go Example

Use direct process launch and debugger-endpoint reuse from Go:

```go
chromeBin := os.Getenv("CPE_BROWSER_CHROME_BIN")
userDataDir := os.Getenv("CPE_BROWSER_AUTOMATION_USER_DATA_DIR")
debugHost := os.Getenv("CPE_BROWSER_REMOTE_DEBUG_HOST")
headlessPort := os.Getenv("CPE_BROWSER_HEADLESS_DEBUG_PORT")
headlessURL := os.Getenv("CPE_BROWSER_HEADLESS_DEBUG_URL")

if chromeBin == "" || userDataDir == "" || debugHost == "" || headlessPort == "" || headlessURL == "" {
    return fmt.Errorf("missing required CPE_BROWSER_* environment variables")
}

client := &http.Client{Timeout: 2 * time.Second}
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, headlessURL, nil)
resp, err := client.Do(req)
if err != nil || resp.StatusCode != http.StatusOK {
    cmd := exec.CommandContext(ctx,
        chromeBin,
        "--user-data-dir="+userDataDir,
        "--remote-debugging-address="+debugHost,
        "--remote-debugging-port="+headlessPort,
        "--headless=new",
    )
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("launch headless Chrome: %w", err)
    }
}
```

For visible mode, switch to `CPE_BROWSER_VISIBLE_DEBUG_PORT` and remove `--headless=new`.

## Operating Rules

- Keep the debugging endpoint bound to `127.0.0.1` only.
- Reuse the same automation `user-data-dir` to preserve site sessions.
- Do not launch visible and headless Chrome simultaneously against the same `user-data-dir`.
- Prefer direct process launch over helper apps or wrapper scripts.
- Default artifacts to `CPE_BROWSER_ARTIFACT_DIR` or the OS temp directory.
- Default viewport to `1600x900` CSS pixels with device scale factor `2` unless the task needs a different breakpoint.

## Practical Note

Chrome 136+ restricts remote debugging against the user's default real profile data directory. Use the separate automation `user-data-dir` instead of the daily-driver Chrome profile for rod/CDP automation.
