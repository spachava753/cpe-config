# Kernel Go SDK Quick Reference

Use this reference when you need the right service or parameter shape quickly without re-reading the whole skill.

## Primary Services

### Client setup

```go
client := kernel.NewClient(option.WithAPIKey(apiKey))
```

The client can also read `KERNEL_API_KEY` from the environment.

### Profiles

```go
profilesPage, err := client.Profiles.List(ctx, kernel.ProfileListParams{
    Limit: kernel.Int(20),
})
```

Use profiles to preload cookies, auth, and browser state.

### Create a browser

```go
browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
    Headless:       kernel.Bool(false),
    TimeoutSeconds: kernel.Int(1800),
    Profile: shared.BrowserProfileParam{
        ID:          kernel.String(profileID),
        SaveChanges: kernel.Bool(false),
    },
})
```

Useful fields on the response:

- `browser.SessionID`
- `browser.BrowserLiveViewURL`
- `browser.CdpWsURL`
- `browser.WebdriverWsURL`

### Reuse an existing browser

```go
info, err := client.Browsers.Get(ctx, sessionID, kernel.BrowserGetParams{})
```

Then continue using `sessionID` with Playwright or computer controls.

### Remote Playwright execution

```go
res, err := client.Browsers.Playwright.Execute(ctx, sessionID, kernel.BrowserPlaywrightExecuteParams{
    Code:       `await page.goto("https://example.com"); return await page.title();`,
    TimeoutSec: kernel.Int(60),
})
```

### Computer controls

```go
_, err := client.Browsers.Computer.CaptureScreenshot(ctx, sessionID, kernel.BrowserComputerCaptureScreenshotParams{})
err = client.Browsers.Computer.ClickMouse(ctx, sessionID, kernel.BrowserComputerClickMouseParams{
    X:      400,
    Y:      250,
    Button: kernel.BrowserComputerClickMouseParamsButtonLeft,
})
err = client.Browsers.Computer.TypeText(ctx, sessionID, kernel.BrowserComputerTypeTextParams{
    Text: "hello world",
})
err = client.Browsers.Computer.PressKey(ctx, sessionID, kernel.BrowserComputerPressKeyParams{
    Key: "Enter",
})
```

### Cleanup

```go
err := client.Browsers.DeleteByID(ctx, sessionID)
```

## Practical Defaults

- Prefer `Headless: kernel.Bool(false)` for visible sessions.
- Prefer `TimeoutSeconds` in the 10-30 minute range for multi-step tasks.
- Prefer `SaveChanges: false` unless persisting profile changes is intentional.
- Prefer `List(...)` over auto-paging when you expect only a few profiles or browsers.
