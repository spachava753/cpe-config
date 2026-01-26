# Gmail Thread Operations

Complete code examples for Gmail thread management using go-imap v1.

## Setup

```go
import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/responses"
)

// Gmail FETCH items - just strings!
const (
	FetchGmailThreadID  imap.FetchItem = "X-GM-THRID"
	FetchGmailMessageID imap.FetchItem = "X-GM-MSGID"
	FetchGmailLabels    imap.FetchItem = "X-GM-LABELS"
)

// Connect helper
func ConnectGmail(email, password string) (*client.Client, error) {
	password = strings.ReplaceAll(password, " ", "")
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

## 1. List Threads

Fetch recent threads from a folder, deduplicate by thread ID.

```go
func ListThreads(c *client.Client, folder string, count int) error {
	// Select folder
	mbox, err := c.Select(folder, true) // true = read-only
	if err != nil {
		return err
	}

	// Determine range
	from := uint32(1)
	if mbox.Messages > uint32(count) {
		from = mbox.Messages - uint32(count) + 1
	}

	// Fetch messages with thread IDs
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, mbox.Messages)

	messages := make(chan *imap.Message, count+10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			FetchGmailThreadID,
			FetchGmailLabels,
		}, messages)
	}()

	// Deduplicate by thread ID
	type ThreadInfo struct {
		ID       string
		Subject  string
		From     string
		MsgCount int
	}
	threads := make(map[string]*ThreadInfo)

	for msg := range messages {
		thrid := fmt.Sprintf("%v", msg.Items[FetchGmailThreadID])
		
		if existing, ok := threads[thrid]; ok {
			existing.MsgCount++
		} else {
			subject := "(no subject)"
			from := ""
			if msg.Envelope != nil {
				subject = msg.Envelope.Subject
				if len(msg.Envelope.From) > 0 {
					from = msg.Envelope.From[0].Address()
				}
			}
			threads[thrid] = &ThreadInfo{
				ID:       thrid,
				Subject:  subject,
				From:     from,
				MsgCount: 1,
			}
		}
	}

	if err := <-done; err != nil {
		return err
	}

	// Display results
	for thrid, info := range threads {
		fmt.Printf("Thread %s: [%d msgs] %s (from: %s)\n",
			thrid, info.MsgCount, info.Subject, info.From)
	}
	return nil
}
```

## 2. Search Threads

Use Gmail's search syntax via X-GM-RAW.

```go
// Custom search command
type GmailRawSearch struct {
	Query string
}

func (s *GmailRawSearch) Command() *imap.Command {
	return &imap.Command{
		Name:      "UID SEARCH",
		Arguments: []interface{}{imap.RawString(`X-GM-RAW "` + s.Query + `"`)},
	}
}

func SearchThreads(c *client.Client, query string) ([]uint32, error) {
	// Select All Mail for comprehensive search
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}

	// Execute search
	searchCmd := &GmailRawSearch{Query: query}
	searchResp := &responses.Search{}
	
	if _, err := c.Execute(searchCmd, searchResp); err != nil {
		return nil, err
	}

	fmt.Printf("Found %d matching messages\n", len(searchResp.Ids))
	return searchResp.Ids, nil
}

// Examples:
// SearchThreads(c, "is:unread")
// SearchThreads(c, "from:github")
// SearchThreads(c, "has:attachment larger:5M")
// SearchThreads(c, "after:2024/01/01 subject:invoice")
```

## 3. Get Thread Contents

Fetch all messages in a thread by thread ID.

```go
// Search by thread ID command
type GmailThreadSearch struct {
	ThreadID string
}

func (s *GmailThreadSearch) Command() *imap.Command {
	return &imap.Command{
		Name:      "UID SEARCH",
		Arguments: []interface{}{imap.RawString("X-GM-THRID " + s.ThreadID)},
	}
}

func GetThread(c *client.Client, threadID string) ([]*imap.Message, error) {
	// Select All Mail to find all messages in thread
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}

	// Search for messages with this thread ID
	searchCmd := &GmailThreadSearch{ThreadID: threadID}
	searchResp := &responses.Search{}
	
	if _, err := c.Execute(searchCmd, searchResp); err != nil {
		return nil, err
	}

	if len(searchResp.Ids) == 0 {
		return nil, fmt.Errorf("no messages in thread %s", threadID)
	}

	// Fetch all messages
	seqset := new(imap.SeqSet)
	for _, uid := range searchResp.Ids {
		seqset.AddNum(uid)
	}

	msgChan := make(chan *imap.Message, len(searchResp.Ids)+1)
	done := make(chan error, 1)
	go func() {
		done <- c.UidFetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			imap.FetchUid,
			FetchGmailLabels,
		}, msgChan)
	}()

	var messages []*imap.Message
	for msg := range msgChan {
		messages = append(messages, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return messages, nil
}
```

## 4. Add Label to Thread

Add a label to all messages in a thread.

```go
// STORE labels command
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
	
	// Quote labels with spaces
	quoted := make([]string, len(s.Labels))
	for i, l := range s.Labels {
		if strings.Contains(l, " ") {
			quoted[i] = `"` + l + `"`
		} else {
			quoted[i] = l
		}
	}
	labelList := "(" + strings.Join(quoted, " ") + ")"
	
	return &imap.Command{
		Name:      "UID STORE",
		Arguments: []interface{}{s.SeqSet, imap.RawString(op), imap.RawString(labelList)},
	}
}

func AddLabelToThread(c *client.Client, threadID string, label string) error {
	// Select in read-write mode
	if _, err := c.Select("[Gmail]/All Mail", false); err != nil {
		return err
	}

	// Find all messages in thread
	searchCmd := &GmailThreadSearch{ThreadID: threadID}
	searchResp := &responses.Search{}
	if _, err := c.Execute(searchCmd, searchResp); err != nil {
		return err
	}

	if len(searchResp.Ids) == 0 {
		return fmt.Errorf("no messages in thread")
	}

	// Build UID set
	seqset := new(imap.SeqSet)
	for _, uid := range searchResp.Ids {
		seqset.AddNum(uid)
	}

	// Add label
	storeCmd := &GmailStoreLabels{
		SeqSet: seqset,
		Add:    true,
		Labels: []string{label},
	}

	if _, err := c.Execute(storeCmd, nil); err != nil {
		return err
	}

	fmt.Printf("Added label '%s' to %d messages\n", label, len(searchResp.Ids))
	return nil
}
```

## 5. Remove Label from Thread

Remove a label from all messages in a thread.

```go
func RemoveLabelFromThread(c *client.Client, threadID string, label string) error {
	// Select in read-write mode
	if _, err := c.Select("[Gmail]/All Mail", false); err != nil {
		return err
	}

	// Find messages
	searchCmd := &GmailThreadSearch{ThreadID: threadID}
	searchResp := &responses.Search{}
	if _, err := c.Execute(searchCmd, searchResp); err != nil {
		return err
	}

	if len(searchResp.Ids) == 0 {
		return fmt.Errorf("no messages in thread")
	}

	// Build UID set
	seqset := new(imap.SeqSet)
	for _, uid := range searchResp.Ids {
		seqset.AddNum(uid)
	}

	// Remove label (Add = false)
	storeCmd := &GmailStoreLabels{
		SeqSet: seqset,
		Add:    false,
		Labels: []string{label},
	}

	if _, err := c.Execute(storeCmd, nil); err != nil {
		return err
	}

	fmt.Printf("Removed label '%s' from %d messages\n", label, len(searchResp.Ids))
	return nil
}
```

## 6. Delete Thread

Move all messages in a thread to Trash.

```go
func DeleteThread(c *client.Client, threadID string) error {
	// Select in read-write mode
	if _, err := c.Select("[Gmail]/All Mail", false); err != nil {
		return err
	}

	// Find messages
	searchCmd := &GmailThreadSearch{ThreadID: threadID}
	searchResp := &responses.Search{}
	if _, err := c.Execute(searchCmd, searchResp); err != nil {
		return err
	}

	if len(searchResp.Ids) == 0 {
		return fmt.Errorf("no messages in thread")
	}

	// Build UID set
	seqset := new(imap.SeqSet)
	for _, uid := range searchResp.Ids {
		seqset.AddNum(uid)
	}

	// Copy to Trash
	if err := c.UidCopy(seqset, "[Gmail]/Trash"); err != nil {
		return err
	}

	// Mark as deleted
	storeItem := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err := c.UidStore(seqset, storeItem, flags, nil); err != nil {
		return err
	}

	// Expunge
	if err := c.Expunge(nil); err != nil {
		return err
	}

	fmt.Printf("Deleted thread %s (%d messages)\n", threadID, len(searchResp.Ids))
	return nil
}
```

## Common Gmail Search Queries

| Query | Description |
|-------|-------------|
| `is:unread` | Unread messages |
| `is:starred` | Starred messages |
| `from:sender@example.com` | From specific sender |
| `to:recipient@example.com` | To specific recipient |
| `subject:keyword` | Subject contains keyword |
| `has:attachment` | Has attachments |
| `larger:5M` | Larger than 5MB |
| `smaller:1M` | Smaller than 1MB |
| `after:2024/01/01` | After date |
| `before:2024/12/31` | Before date |
| `label:Work` | Has specific label |
| `in:inbox` | In inbox |
| `in:sent` | In sent mail |
