# Gmail SMTP Sending Operations

Complete code examples for sending emails via Gmail SMTP using `emersion/go-smtp` and `emersion/go-sasl`.

For end-to-end thread-aware replies that start with IMAP discovery and then send over SMTP, see `replying.md`.

## Setup

```go
import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

const (
	GmailSMTPHost = "smtp.gmail.com"
	GmailSMTPPort = "465" // SSL/TLS
)

// ConnectSMTP establishes an authenticated SMTP connection to Gmail
func ConnectSMTP(email, appPassword string) (*smtp.Client, error) {
	// Remove spaces from app password (Google formats it with spaces)
	appPassword = strings.ReplaceAll(appPassword, " ", "")

	// Connect with TLS directly on port 465
	tlsConfig := &tls.Config{ServerName: GmailSMTPHost}
	conn, err := tls.Dial("tcp", GmailSMTPHost+":"+GmailSMTPPort, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("TLS dial failed: %w", err)
	}

	// Create SMTP client from TLS connection
	c := smtp.NewClient(conn)

	// Authenticate
	auth := sasl.NewPlainClient("", email, appPassword)
	if err := c.Auth(auth); err != nil {
		c.Close()
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	return c, nil
}
```

## 1. Send Plain Text Email

Simple text email without attachments.

```go
func SendTextEmail(from, to, subject, body string, appPassword string) error {
	c, err := ConnectSMTP(from, appPassword)
	if err != nil {
		return err
	}
	defer c.Close()

	// Build RFC 5322 message
	msg := buildMessage(from, to, subject, "text/plain", body)

	// Send
	if err := c.Mail(from, nil); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	if err := c.Rcpt(to, nil); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("close failed: %w", err)
	}

	return c.Quit()
}

// buildMessage creates an RFC 5322 compliant email message
func buildMessage(from, to, subject, contentType, body string) string {
	// Generate Message-ID
	msgID := fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), from)

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
		fmt.Sprintf("Message-ID: %s", msgID),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: %s; charset=UTF-8", contentType),
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
```

## 2. Send HTML Email

Email with HTML content.

```go
func SendHTMLEmail(from, to, subject, htmlBody string, appPassword string) error {
	c, err := ConnectSMTP(from, appPassword)
	if err != nil {
		return err
	}
	defer c.Close()

	msg := buildMessage(from, to, subject, "text/html", htmlBody)

	if err := c.Mail(from, nil); err != nil {
		return err
	}
	if err := c.Rcpt(to, nil); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

// Example usage:
// html := "<h1>Hello!</h1><p>This is an <strong>HTML</strong> email.</p>"
// SendHTMLEmail("me@gmail.com", "you@example.com", "Test HTML", html, appPassword)
```

## 3. Send to Multiple Recipients

Email with To, CC, and BCC recipients.

```go
func SendToMultiple(from string, to, cc, bcc []string, subject, body string, appPassword string) error {
	c, err := ConnectSMTP(from, appPassword)
	if err != nil {
		return err
	}
	defer c.Close()

	// Build message with CC (BCC is not included in headers)
	msg := buildMultiRecipientMessage(from, to, cc, subject, body)

	if err := c.Mail(from, nil); err != nil {
		return err
	}

	// Add all recipients (To, CC, and BCC)
	allRecipients := append(append(to, cc...), bcc...)
	for _, rcpt := range allRecipients {
		if err := c.Rcpt(rcpt, nil); err != nil {
			return fmt.Errorf("RCPT TO %s failed: %w", rcpt, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

func buildMultiRecipientMessage(from string, to, cc []string, subject, body string) string {
	msgID := fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), from)

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", strings.Join(to, ", ")),
	}

	if len(cc) > 0 {
		headers = append(headers, fmt.Sprintf("Cc: %s", strings.Join(cc, ", ")))
	}

	headers = append(headers,
		fmt.Sprintf("Subject: %s", subject),
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
		fmt.Sprintf("Message-ID: %s", msgID),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	)

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
```

## 4. Send Email with Attachment

Email with file attachments using MIME multipart.

```go
import (
	"encoding/base64"
	"mime"
	"os"
	"path/filepath"
)

func SendWithAttachment(from, to, subject, body string, attachmentPaths []string, appPassword string) error {
	c, err := ConnectSMTP(from, appPassword)
	if err != nil {
		return err
	}
	defer c.Close()

	msg, err := buildMultipartMessage(from, to, subject, body, attachmentPaths)
	if err != nil {
		return err
	}

	if err := c.Mail(from, nil); err != nil {
		return err
	}
	if err := c.Rcpt(to, nil); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

func buildMultipartMessage(from, to, subject, body string, attachments []string) (string, error) {
	boundary := fmt.Sprintf("boundary_%d", time.Now().UnixNano())
	msgID := fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), from)

	var sb strings.Builder

	// Headers
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	sb.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	sb.WriteString(fmt.Sprintf("Message-ID: %s\r\n", msgID))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	sb.WriteString("\r\n")

	// Body part
	sb.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	sb.WriteString("\r\n")

	// Attachments
	for _, path := range attachments {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read attachment %s: %w", path, err)
		}

		filename := filepath.Base(path)
		mimeType := mime.TypeByExtension(filepath.Ext(path))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		sb.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		sb.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", mimeType, filename))
		sb.WriteString("Content-Transfer-Encoding: base64\r\n")
		sb.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename))
		sb.WriteString("\r\n")

		// Base64 encode with line wrapping (76 chars per line)
		encoded := base64.StdEncoding.EncodeToString(data)
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			sb.WriteString(encoded[i:end])
			sb.WriteString("\r\n")
		}
	}

	// Final boundary
	sb.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return sb.String(), nil
}
```

## 5. Reply to Email

Reply to an existing email (In-Reply-To and References headers).

This section shows the SMTP half only. For a reliable Gmail reply in an existing conversation, first fetch the target message's `Message-ID` and existing `References` over IMAP, then append that chain when building the outgoing message. See `replying.md`.

```go
func SendReply(from, to, subject, body, originalMessageID, existingReferences string, appPassword string) error {
	c, err := ConnectSMTP(from, appPassword)
	if err != nil {
		return err
	}
	defer c.Close()

	msg := buildReplyMessage(from, to, subject, body, originalMessageID, existingReferences)

	if err := c.Mail(from, nil); err != nil {
		return err
	}
	if err := c.Rcpt(to, nil); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

func buildReplyMessage(from, to, subject, body, originalMsgID, existingReferences string) string {
	msgID := fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), from)

	// Ensure subject has "Re:" prefix
	if !strings.HasPrefix(strings.ToLower(subject), "re:") {
		subject = "Re: " + subject
	}

	references := strings.TrimSpace(strings.TrimSpace(existingReferences + " " + originalMsgID))
	references = strings.Join(strings.Fields(references), " ")

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
		fmt.Sprintf("Message-ID: %s", msgID),
		fmt.Sprintf("In-Reply-To: %s", originalMsgID),
		fmt.Sprintf("References: %s", references),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
```

## 6. Complete Example

Full working example with error handling.

```go
package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

func main() {
	email := os.Getenv("GMAIL_ADDRESS")       // or hardcode
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")

	if email == "" || appPassword == "" {
		fmt.Println("Set GMAIL_ADDRESS and GMAIL_APP_PASSWORD")
		os.Exit(1)
	}

	err := SendTextEmail(
		email,
		"recipient@example.com",
		"Hello from Go!",
		"This is a test email sent using emersion/go-smtp.",
		appPassword,
	)

	if err != nil {
		fmt.Printf("Failed to send: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Email sent successfully!")
}

// ... include ConnectSMTP, SendTextEmail, buildMessage functions from above
```

## Gmail SMTP Limits

| Account Type | Daily Limit |
|--------------|-------------|
| Free Gmail (@gmail.com) | 500 emails/day |
| Google Workspace | 2,000 emails/day |

## Troubleshooting

| Error | Solution |
|-------|----------|
| `auth failed` | Check app password (no spaces), ensure 2FA is enabled |
| `connection refused` | Check firewall, try port 465 with TLS |
| `certificate error` | Ensure correct ServerName in TLS config |
| `assignment mismatch` or `too many arguments` around `smtp.NewClient` | You may be mixing `emersion/go-smtp` with `net/smtp` APIs; in `emersion/go-smtp`, `smtp.NewClient` takes only the connection |
| `too many login attempts` | Wait and retry, check for account security alerts |

## Dependencies

Add to your `go.mod`:

```
require (
	github.com/emersion/go-sasl v0.0.0-20231106173351-e73c9f7bad43
	github.com/emersion/go-smtp v0.21.3
)
```

Or run:

```bash
go get github.com/emersion/go-sasl
go get github.com/emersion/go-smtp
```
