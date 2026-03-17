---
name: kernel-browser-use
description: Use when tasks require creating or reusing Kernel remote browser sessions from Go, driving pages with the Kernel Go SDK, or falling back to computer-use controls. Covers client setup, profile selection, browser lifecycle, Playwright executor usage, visible live-view sessions, and practical SDK caveats.
---

# Kernel Browser Use

Use this skill for general browser automation and browser-assisted investigation with Kernel. Prefer it when a task benefits from a remote browser with reusable authenticated state, a visible live view, or a mix of scripted automation and direct computer-use controls.

## Execution Posture

- Default posture: create one visible browser session, keep it alive across related work, and reuse its `session_id` instead of creating and deleting a new browser in every step.
- Browser creation phase: prefer `Headless: kernel.Bool(false)` so the user can watch the browser in Kernel live view or the dashboard. Use a generous `TimeoutSeconds`.
- Main automation phase: prefer `client.Browsers.Playwright.Execute(...)` for navigation, inspection, extraction, and authenticated app work. Return structured data instead of large raw dumps.
- Fallback phase: use `client.Browsers.Computer` controls when Playwright is not enough, especially for canvas-heavy UIs, drag-and-drop, coordinate-sensitive flows, or native-feeling interactions.
- Cleanup phase: delete the browser only when the task is actually done, the user asks, or the session is clearly disposable.

## Recommended Workflow

1. Ensure `KERNEL_API_KEY` is present in the environment.
2. Create a Kernel client with `kernel.NewClient(...)`.
3. List profiles and select the intended profile by ID or name.
4. Create a browser with that profile, `Headless: false`, and a timeout long enough for the whole task.
5. Record and surface the returned `session_id` and `browser_live_view_url`.
6. Reuse that same `session_id` in later `execute_go_code` calls.
7. Drive the browser with `Playwright.Execute` whenever possible.
8. Use screenshots and live view to validate progress after important transitions.
9. Clean up the session at the end if appropriate.

## Core Rules

### Create once, reuse often

- Do not recreate or delete the browser between every tool call.
- Keep the `session_id` in working context and pass it into later SDK calls.
- If needed, verify the session still exists with `client.Browsers.Get(ctx, sessionID, kernel.BrowserGetParams{})`.

### Prefer the Playwright executor over raw CDP for most work

- The SDK's Playwright executor runs next to the browser and already gives access to `page`, `context`, and `browser`.
- This is usually simpler and more reliable than wiring up a separate local Playwright/CDP stack.
- Use raw CDP only when the task explicitly needs DevTools protocol features or a local protocol client.

### Keep sessions visible by default

- Prefer `Headless: kernel.Bool(false)` so the user can see what is happening.
- Always surface `BrowserLiveViewURL` after browser creation when available.
- Only switch to `Headless: true` when the task specifically benefits from it.

### Return structured data from Playwright scripts

- Prefer small JSON-like return values instead of large text blobs.
- Use locator APIs for page text when possible.
- For DOM access, use `page.evaluate(...)` or Playwright locators rather than assuming page globals exist in the script body.

### Use app data when it is more reliable than rendered text

- On authenticated SPAs, it can be more robust to inspect fetch/XHR responses or call the app's own API from inside the browser context.
- Use visible scraping only when the rendered UI is the real source of truth or when app data is inaccessible.

## When to Use Computer Controls

Use `client.Browsers.Computer` when:

- a flow requires mouse coordinates or drag-and-drop
- the page is canvas/WebGL-heavy and hard to address with locators
- a site blocks scripted interactions but still responds to ordinary pointer and keyboard events
- visual confirmation through screenshots is important

Typical loop:

1. Capture a screenshot.
2. Decide on coordinates or a scroll target.
3. Click, type, scroll, or press keys.
4. Capture another screenshot or inspect state with Playwright.

## Important SDK Caveats

- `Playwright.Execute` does not run as raw page JavaScript. Its code runs in a function with `page`, `context`, and `browser`. Top-level DOM globals like `document` are not available unless you access them through `page.evaluate(...)`.
- The SDK's auto-pager can fail if the server omits or leaves `X-Next-Offset` empty. For small lists, prefer `List(...)` and inspect `page.Items` directly.
- Browser profile selection uses `shared.BrowserProfileParam`; provide either `ID` or `Name`.
- If profile state should persist back to Kernel after the session ends, set `SaveChanges: kernel.Bool(true)` explicitly. Otherwise, default to `false`.
- Because each `execute_go_code` call runs in a fresh process, persistence across calls means reusing the `session_id`, not keeping a Go client object alive.

## References

- `references/api_reference.md` - service and type quick reference
- `references/go-sdk-examples.md` - reusable code examples
- `references/sdk-quirks.md` - caveats and practical workarounds
- `scripts/example.go` - runnable Go example for creating and reusing a browser session
