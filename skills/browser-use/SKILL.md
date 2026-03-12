---
name: browser-use
description: Use when tasks require driving Chrome with rod/CDP for navigation, scraping, frontend QA, visual verification, screenshots, or proof-of-work recordings. Default to a dedicated automation profile and headless mode, prefer rod APIs over OS keystrokes or app activation, and hand control to the user for logins, 1Password, MFA, captchas, device verification, passkeys, or other interactive auth steps.
---

# Browser Use

## Overview

Drive Chrome through `github.com/go-rod/rod` using a dedicated automation profile, with headless mode as the default. Treat rendered-page inspection as first-class: prefer screenshots and recordings of the page over DOM-only inspection, and avoid foreground focus steals unless the user explicitly asks for a visible browser window.

## Environment Configuration

Do not hard-code machine-specific absolute paths in the skill. Read browser settings from environment variables instead:

- `CPE_BROWSER_CHROME_BIN`: Chrome binary to launch
- `CPE_BROWSER_AUTOMATION_USER_DATA_DIR`: dedicated automation profile directory
- `CPE_BROWSER_REMOTE_DEBUG_HOST`: host for remote debugging, normally `127.0.0.1`
- `CPE_BROWSER_VISIBLE_DEBUG_PORT`: port for visible Chrome automation
- `CPE_BROWSER_HEADLESS_DEBUG_PORT`: port for headless Chrome automation
- `CPE_BROWSER_VISIBLE_DEBUG_URL`: browser version endpoint for visible Chrome
- `CPE_BROWSER_HEADLESS_DEBUG_URL`: browser version endpoint for headless Chrome
- `CPE_BROWSER_ARTIFACT_DIR`: default artifact directory; prefer this or `os.TempDir()` unless the user requests a specific folder
- `CPE_BROWSER_VIEWPORT_WIDTH`: default CSS viewport width for browser captures
- `CPE_BROWSER_VIEWPORT_HEIGHT`: default CSS viewport height for browser captures
- `CPE_BROWSER_DEVICE_SCALE_FACTOR`: default pixel density for screenshots and recordings

Fail early with a clear error if the required browser environment variables are missing. Launch Chrome directly from the agent as needed. Do not rely on helper launcher apps or wrapper scripts unless the user explicitly asks for them.

Do not run two Chrome processes against the same `--user-data-dir` at the same time. If switching between visible and headless modes, stop one before starting the other.

## Default Policy

1. Default to headless mode.
2. Use visible mode only when the user explicitly requests it, or when the workflow requires a human-visible browser window.
3. Prefer rod/CDP actions such as navigation, clicking, typing, waiting, screenshots, and downloads.
4. Avoid AppleScript, `System Events`, global keystrokes, app activation, and other focus-stealing OS automation unless the user explicitly approves it for that task.
5. Launch Chrome directly with `exec.Command` and the configured environment variables; no persistent launcher helper is required.
6. Reuse the same automation `user-data-dir` so session cookies and site logins persist across runs.
7. Before launching Chrome, check whether the configured visible or headless debugger endpoint is already live and reuse it when possible.
8. Produce proof-of-work artifacts for important tasks or whenever the user requests them.

## Mode Selection

Choose the mode deliberately:

- Use headless for routine browsing, scraping, frontend QA, screenshot capture, visual comparison, and multi-step navigation that does not need extension UI or native prompts.
- Use visible mode when the user wants to watch, when browser extensions or native dialogs must be interacted with, or when a site presents login, MFA, CAPTCHA, device verification, passkey, or similar handoff steps.
- If the task is ambiguous, start headless and switch to visible only if a blocking interactive step appears.

## Recording And Proof Of Work

When proof of work matters, capture rendered output rather than relying only on DOM inspection.

- Default artifact location is `CPE_BROWSER_ARTIFACT_DIR` or the OS temp directory unless the user explicitly asks for a specific folder.
- Name artifacts with a semantic task slug plus a Unix timestamp, such as `<task-name>-<unix-timestamp>.gif`, `<task-name>-<unix-timestamp>.mp4`, or `<task-name>-<unix-timestamp>-0001.png`.
- For short tasks, capture step screenshots and assemble a GIF if helpful.
- For longer tasks, prefer MP4 when a suitable encoder is available; otherwise use a GIF or a directory of ordered frames.
- Capture key transition points even when a full recording is not needed: initial page, after each significant action, and final state.
- Treat screenshots and recordings as potentially sensitive because they may include account data, messages, or secrets visible in the page.
- Mask or avoid recording sensitive form entry when possible.
- Unless the user asks to keep artifacts, prefer temporary storage and clean up no-longer-needed intermediates.

## Authentication And Handoff

Do not try to force interactive authentication through background tricks.

- If the page requires login, 1Password interaction, MFA, email code entry, CAPTCHA, device/browser verification, or other human approval, stop and hand control to the user.
- Explain what blocked progress and what the user should do next.
- Once the user finishes the interactive step in the automation profile, reconnect with rod and continue.

Headless-to-visible switching is not hot-swappable in a running Chrome process.

- A running headless Chrome process cannot simply be turned into a visible window.
- To hand off from headless to visible, save the current URL and a screenshot, stop the headless process, relaunch Chrome in visible mode with the same `--user-data-dir`, and reopen the relevant URL.
- Persisted cookies, local storage, and logged-in site sessions normally survive because the same `user-data-dir` is reused.
- Transient in-memory page state may not survive. If an exact live flow is critical, prefer visible mode from the outset.

## Browser Workflow

Follow this workflow unless the task requires a variation:

1. Validate required environment variables before launch.
2. Determine whether the task can stay headless or is likely to need visible handoff.
3. Check whether the chosen debugger endpoint is already live; reuse it when possible, otherwise launch Chrome in the selected mode.
4. Apply the configured viewport and device scale defaults so captures stay consistent across runs.
5. Open a dedicated tab or target for the task rather than disturbing unrelated tabs.
6. Navigate with rod/CDP only; avoid frontmost-app shortcuts.
7. Wait for stable page state, inspect rendered output with screenshots when needed, and collect only the evidence needed for the task.
8. If an interactive auth or verification step appears, pause and hand back to the user.
9. Resume after handoff, finish the task, and return concise findings plus any requested artifacts.

## References

Use these bundled references when needed:

- `references/modes-and-recording.md` for choosing headless vs visible mode and selecting proof-of-work artifacts.
- `references/auth-and-handoff.md` for login blockers and the exact headless-to-visible handoff pattern.
- `references/local-setup.md` for the environment variables, launch patterns, debugger endpoint conventions, and viewport defaults.
