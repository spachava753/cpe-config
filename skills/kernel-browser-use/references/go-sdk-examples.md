# Kernel Go SDK Examples

These examples are intentionally general. Adapt the target URL, selectors, profile selection, and data extraction to the task at hand.

Assume `KERNEL_API_KEY` is already set in the environment.

## 1. Create a visible browser from a profile

```go
package main

import (
    "context"
    "fmt"

    kernel "github.com/kernel/kernel-go-sdk"
    "github.com/kernel/kernel-go-sdk/option"
    "github.com/kernel/kernel-go-sdk/shared"
)

func main() {
    ctx := context.Background()
    client := kernel.NewClient(option.WithAPIKey("YOUR_API_KEY"))

    profilesPage, err := client.Profiles.List(ctx, kernel.ProfileListParams{Limit: kernel.Int(20)})
    if err != nil {
        panic(err)
    }
    if len(profilesPage.Items) == 0 {
        panic("no profiles found")
    }

    profile := profilesPage.Items[0]
    browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
        Headless:       kernel.Bool(false),
        TimeoutSeconds: kernel.Int(1800),
        Profile: shared.BrowserProfileParam{
            ID:          kernel.String(profile.ID),
            SaveChanges: kernel.Bool(false),
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("session:", browser.SessionID)
    fmt.Println("live view:", browser.BrowserLiveViewURL)
}
```

## 2. Reuse an existing session in a later step

```go
info, err := client.Browsers.Get(ctx, sessionID, kernel.BrowserGetParams{})
if err != nil {
    panic(err)
}
fmt.Println(info.SessionID, info.BrowserLiveViewURL)
```

Once the session exists, use the same `sessionID` in later calls rather than creating a new browser.

## 3. Run a small Playwright script

```go
res, err := client.Browsers.Playwright.Execute(ctx, sessionID, kernel.BrowserPlaywrightExecuteParams{
    TimeoutSec: kernel.Int(90),
    Code: `
        await page.goto('https://example.com', { waitUntil: 'domcontentloaded' });
        return {
          url: page.url(),
          title: await page.title(),
          heading: await page.locator('h1').first().innerText().catch(() => '')
        };
    `,
})
if err != nil {
    panic(err)
}
fmt.Printf("success=%v result=%#v\n", res.Success, res.Result)
```

## 4. Access DOM state correctly from Playwright executor code

Use `page.evaluate(...)` or Playwright locators. Do not assume top-level `document` is available.

```go
res, err := client.Browsers.Playwright.Execute(ctx, sessionID, kernel.BrowserPlaywrightExecuteParams{
    TimeoutSec: kernel.Int(90),
    Code: `
        await page.goto('https://example.com');
        const title = await page.evaluate(() => document.title);
        const bodyText = await page.locator('body').innerText();
        return { title, bodyPreview: bodyText.slice(0, 300) };
    `,
})
```

## 5. Inspect authenticated app data instead of scraping visible text

For SPAs, it is often more reliable to query app data from inside the session.

```go
res, err := client.Browsers.Playwright.Execute(ctx, sessionID, kernel.BrowserPlaywrightExecuteParams{
    TimeoutSec: kernel.Int(120),
    Code: `
        await page.goto('https://app.example.com', { waitUntil: 'domcontentloaded' });
        const data = await page.evaluate(async () => {
          const resp = await fetch('/api/me', { credentials: 'include' });
          return await resp.json();
        });
        return data;
    `,
})
```

## 6. Watch network responses during a page action

```go
res, err := client.Browsers.Playwright.Execute(ctx, sessionID, kernel.BrowserPlaywrightExecuteParams{
    TimeoutSec: kernel.Int(120),
    Code: `
        const seen = [];
        page.on('response', async (resp) => {
          const url = resp.url();
          const ct = (resp.headers()['content-type'] || '').toLowerCase();
          if (!ct.includes('json')) return;
          if (!/api|graphql/.test(url)) return;
          const body = await resp.text().catch(() => '');
          seen.push({ url, status: resp.status(), bodyPreview: body.slice(0, 500) });
        });
        await page.goto('https://app.example.com');
        await page.waitForTimeout(5000);
        return seen;
    `,
})
```

## 7. Use computer controls as a fallback

```go
shot, err := client.Browsers.Computer.CaptureScreenshot(ctx, sessionID, kernel.BrowserComputerCaptureScreenshotParams{})
if err != nil {
    panic(err)
}
_ = shot

err = client.Browsers.Computer.ClickMouse(ctx, sessionID, kernel.BrowserComputerClickMouseParams{
    X:      640,
    Y:      380,
    Button: kernel.BrowserComputerClickMouseParamsButtonLeft,
})
if err != nil {
    panic(err)
}

err = client.Browsers.Computer.TypeText(ctx, sessionID, kernel.BrowserComputerTypeTextParams{Text: "search query"})
if err != nil {
    panic(err)
}

err = client.Browsers.Computer.PressKey(ctx, sessionID, kernel.BrowserComputerPressKeyParams{Key: "Enter"})
if err != nil {
    panic(err)
}
```

## 8. Clean up only at the end

```go
if done {
    if err := client.Browsers.DeleteByID(ctx, sessionID); err != nil {
        panic(err)
    }
}
```

For multi-step work, avoid `defer client.Browsers.DeleteByID(...)` in the same code path that creates the session unless the whole task is meant to be one-shot.
