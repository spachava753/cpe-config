# Gmail IMAP API Reference

Gmail supports normal IMAP operations. This file focuses on the Gmail-specific extensions exposed through IMAP. For behavior notes and practical caveats, also see `gmail-extensions.md`.

## go-imap v1 Gmail Extensions

This reference shows how to use Gmail IMAP extensions with `github.com/emersion/go-imap` v1.

### Custom Fetch Items

```go
import "github.com/emersion/go-imap"

// Gmail FETCH items - just strings that go-imap stores automatically
const (
    FetchGmailThreadID  imap.FetchItem = "X-GM-THRID"
    FetchGmailMessageID imap.FetchItem = "X-GM-MSGID"
    FetchGmailLabels    imap.FetchItem = "X-GM-LABELS"
)

// Usage in Fetch
items := []imap.FetchItem{imap.FetchEnvelope, FetchGmailThreadID, FetchGmailLabels}
c.Fetch(seqset, items, messages)

// Access results - automatically parsed into msg.Items map
for msg := range messages {
    threadID := msg.Items[FetchGmailThreadID]  // interface{}
    labels := msg.Items[FetchGmailLabels]      // interface{} (usually []interface{})
}
```

### Custom Search Commands

```go
import (
    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
    "github.com/emersion/go-imap/responses"
)

// Search by thread ID
type GmailThreadSearch struct {
    ThreadID string
}

func (s *GmailThreadSearch) Command() *imap.Command {
    return &imap.Command{
        Name:      "UID SEARCH",
        Arguments: []interface{}{imap.RawString("X-GM-THRID " + s.ThreadID)},
    }
}

// Search with Gmail syntax
type GmailRawSearch struct {
    Query string
}

func (s *GmailRawSearch) Command() *imap.Command {
    return &imap.Command{
        Name:      "UID SEARCH",
        Arguments: []interface{}{imap.RawString(`X-GM-RAW "` + s.Query + `"`)},
    }
}

// Execute search
searchCmd := &GmailThreadSearch{ThreadID: "1234567890"}
searchResp := &responses.Search{}
status, err := c.Execute(searchCmd, searchResp)
// searchResp.Ids contains matching UIDs
```

### Custom Store Commands (Labels)

```go
// Add/remove labels
type GmailStoreLabels struct {
    SeqSet *imap.SeqSet
    Add    bool     // true = +X-GM-LABELS, false = -X-GM-LABELS
    Labels []string
}

func (s *GmailStoreLabels) Command() *imap.Command {
    op := "-X-GM-LABELS"
    if s.Add {
        op = "+X-GM-LABELS"
    }
    labelList := "(" + strings.Join(s.Labels, " ") + ")"
    return &imap.Command{
        Name:      "UID STORE",
        Arguments: []interface{}{s.SeqSet, imap.RawString(op), imap.RawString(labelList)},
    }
}
```

### Connection Helper

```go
func ConnectGmail(email, password string) (*client.Client, error) {
    c, err := client.DialTLS("imap.gmail.com:993", nil)
    if err != nil {
        return nil, err
    }
    if err := c.Login(email, password); err != nil {
        c.Logout()
        return nil, err
    }
    return c, nil
}
```

## Gmail Special Folders

| Folder | Purpose |
|--------|---------|
| `INBOX` | Main inbox |
| `[Gmail]/All Mail` | All messages |
| `[Gmail]/Sent Mail` | Sent messages |
| `[Gmail]/Drafts` | Drafts |
| `[Gmail]/Trash` | Deleted messages |
| `[Gmail]/Spam` | Spam |
| `[Gmail]/Starred` | Starred messages |
| `[Gmail]/Important` | Important messages |
