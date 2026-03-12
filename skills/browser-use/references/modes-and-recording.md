# Modes And Recording

## Mode Decision Rules

Prefer the least disruptive mode that can still complete the task.

Choose `headless` when:
- The task is routine browsing, scraping, clicking, form filling without interactive auth, screenshot capture, or rendered-page verification.
- The user does not need to watch the browser live.
- Avoiding foreground focus steals is important.

Choose `visible` when:
- The user explicitly asks to watch the browser.
- The site or workflow needs extension UI, native file pickers, permission prompts, passkeys, or other browser chrome interactions.
- The task hits login, MFA, CAPTCHA, device verification, or any step where the user must take over.

## Proof Of Work

Rendered-page proof is preferred over DOM-only claims.

Recommended artifact policy:
- `single screenshot`: quick verification of a final state.
- `step screenshots`: before/after key actions for short or medium tasks.
- `GIF`: short action sequences where motion matters and file size should stay modest.
- `MP4`: longer sessions when smoother playback matters and an encoder is available.
- `ordered frames`: fallback when no encoder is available or when the user may want frame-by-frame review.

Default storage and naming:
- Put artifacts in `CPE_BROWSER_ARTIFACT_DIR` or the OS temp directory unless the user asks for a specific destination.
- Use semantic task slugs plus Unix timestamps, such as `<task-name>-<unix-timestamp>.gif`.
- Prefer cleaning up intermediate frame files after producing the final requested artifact.

## Capture Strategy

Use a hybrid capture strategy by default:
- Capture an initial frame.
- Capture after each significant action.
- Capture on errors or unexpected UI.
- Capture the final frame.
- Add timed frames only when a smoother replay is useful.

## Visual Consistency

Use a consistent viewport so screenshots remain comparable across tasks.

- Default to `1600x900` CSS pixels with device scale factor `2`.
- This yields high-fidelity captures while preserving a standard desktop layout and readable text.
- For downstream model analysis on `gpt-5.4`, prefer `detail: "original"` when available.
- If a task requires more vertical context, prefer scrolling plus multiple captures over making the viewport extremely tall.
- If a site only reproduces a bug at a specific breakpoint, match that breakpoint intentionally and report it.

## Rendered-Page Caveat

Headless captures show the rendered page content, not the full desktop or every browser chrome surface.

- Good fit: page layout, text, images, modals, menus rendered in-page, scrolling, and visual regressions.
- Less reliable: extension popovers, native file pickers, OS permission prompts, browser chrome UI, and some system-auth surfaces.

## Safety Notes

Screenshots and recordings can expose sensitive information.

- Avoid recording password entry if possible.
- Prefer pausing recording during sensitive input.
- Mention when the recording may contain account-specific data.
- Keep artifacts scoped to the task rather than capturing unnecessary browsing.
- Do not print cookies, tokens, local-storage secrets, autofilled credentials, or raw auth headers into logs or normal task output.
