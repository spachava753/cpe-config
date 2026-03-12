# Authentication And Handoff

## When To Stop And Hand Back

Stop automation and ask the user to take over when any of these appear:
- Login pages that depend on password managers, extensions, passkeys, or manual approval
- MFA prompts, security keys, email or SMS verification, device verification, or "verify this browser"
- CAPTCHAs or human checks
- Native dialogs or browser UI that rod/CDP cannot reliably control in headless mode

## Headless To Visible Transition

A running headless Chrome process cannot be converted into a visible one without restarting the process.

Use this handoff sequence:
1. Record the current URL.
2. Capture a screenshot of the blocked state.
3. Stop the headless browser process.
4. Launch visible Chrome with the same `--user-data-dir`.
5. Reopen the same URL or the nearest stable page in the flow.
6. Let the user complete the interactive step.
7. Reconnect with rod and continue.

## What Persists Across The Restart

When the same automation `user-data-dir` is reused, the following usually persist:
- Site cookies and signed-in sessions
- Local storage and extension state
- Browser profile settings

The following may not persist reliably:
- Exact transient in-memory JavaScript state
- Unsaved form data
- A challenge page that depends on a live request flow

If the interactive step is likely to depend on transient state, start visible from the beginning instead of switching later.

## Communication Pattern

When pausing for handoff, report:
- What blocked progress
- Whether the browser should remain visible for the next steps
- Exactly what the user should do
- How the user should tell the agent to resume
