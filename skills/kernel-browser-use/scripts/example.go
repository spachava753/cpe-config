//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/shared"
)

// Example: create a visible Kernel browser session from a profile and run a
// small Playwright script against it.
//
// Required environment variables:
//   KERNEL_API_KEY
//
// Optional environment variables:
//   KERNEL_PROFILE_NAME   - profile name to match; otherwise uses the first profile
//   KERNEL_KEEP_OPEN=1    - keep the browser session open after the script exits
//
// This example intentionally prefers a visible browser and leaves the session
// open by default so it can be reused.

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("KERNEL_API_KEY")
	if apiKey == "" {
		panic("KERNEL_API_KEY is required")
	}

	client := kernel.NewClient(option.WithAPIKey(apiKey))

	profilesPage, err := client.Profiles.List(ctx, kernel.ProfileListParams{Limit: kernel.Int(20)})
	if err != nil {
		panic(err)
	}
	if len(profilesPage.Items) == 0 {
		panic("no Kernel profiles found")
	}

	profileName := os.Getenv("KERNEL_PROFILE_NAME")
	profileID := profilesPage.Items[0].ID
	profileLabel := profilesPage.Items[0].Name
	if profileName != "" {
		found := false
		for _, p := range profilesPage.Items {
			if p.Name == profileName {
				profileID = p.ID
				profileLabel = p.Name
				found = true
				break
			}
		}
		if !found {
			panic("requested profile name was not found")
		}
	}

	browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
		Headless:       kernel.Bool(false),
		TimeoutSeconds: kernel.Int(1800),
		Profile: shared.BrowserProfileParam{
			ID:          kernel.String(profileID),
			SaveChanges: kernel.Bool(false),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("session_id=%s\n", browser.SessionID)
	fmt.Printf("live_view=%s\n", browser.BrowserLiveViewURL)
	fmt.Printf("profile=%q\n", profileLabel)

	result, err := client.Browsers.Playwright.Execute(ctx, browser.SessionID, kernel.BrowserPlaywrightExecuteParams{
		TimeoutSec: kernel.Int(90),
		Code: `
			await page.goto('https://example.com', { waitUntil: 'domcontentloaded' });
			return {
			  url: page.url(),
			  title: await page.title(),
			  bodyPreview: (await page.locator('body').innerText()).slice(0, 200)
			};
		`,
	})
	if err != nil {
		panic(err)
	}

	encoded, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(encoded))

	if os.Getenv("KERNEL_KEEP_OPEN") == "1" || os.Getenv("KERNEL_KEEP_OPEN") == "" {
		fmt.Println("browser left open for reuse")
		return
	}

	if err := client.Browsers.DeleteByID(ctx, browser.SessionID); err != nil {
		panic(err)
	}
	fmt.Println("browser deleted")
}
