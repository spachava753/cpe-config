# Kernel SDK Quirks and Workarounds

## 1. Prefer visible sessions unless a task explicitly needs headless

A visible browser is usually more useful for CPE work because the user can inspect the live view and confirm progress. Use `Headless: kernel.Bool(false)` as the default choice.

## 2. Reuse the browser session across tool calls

Each `execute_go_code` call runs in a fresh Go process. If you want continuity, keep the browser alive on the Kernel side and reuse the `session_id` in later calls.

Practical pattern:

- create browser once
- print `session_id` and `browser_live_view_url`
- store those values in working context
- use the same `session_id` in later Playwright or computer-control calls
- delete only when the task is complete

## 3. Playwright executor code is not raw page JavaScript

The executor gives you `page`, `context`, and `browser`. Top-level DOM globals are not directly in scope.

Bad pattern:

```js
const title = document.title
```

Good patterns:

```js
const title = await page.evaluate(() => document.title)
```

```js
const bodyText = await page.locator('body').innerText()
```

## 4. Auto-paging can be brittle

The SDK auto-pager relies on `X-Next-Offset`. If the service returns an empty or missing header, auto-paging can fail. For small profile or browser lists, prefer a single `List(...)` call and inspect `.Items` directly.

## 5. Prefer structured returns from Playwright

Do not return whole documents, giant arrays, or large page text unless needed. Return a small object with the exact fields required for the next step.

Example:

```js
return {
  url: page.url(),
  title: await page.title(),
  items: extracted.slice(0, 20)
}
```

## 6. Use authenticated app APIs when the UI is noisy

Some applications render dense UIs where visible text is incomplete or repetitive. If the task permits it, inspect app fetch/XHR traffic or call the app's API from inside the browser session.

## 7. Use computer controls for interaction, not primary extraction

Computer controls are best for clicking, dragging, typing, or validating visuals. For extraction, Playwright and app-data access are usually more reliable and less coordinate-sensitive.

## 8. Be explicit about profile persistence

If the task should update cookies, storage, or other state back into the Kernel profile, set:

```go
SaveChanges: kernel.Bool(true)
```

Otherwise, prefer `false` so exploratory browsing does not unintentionally modify the saved profile.
