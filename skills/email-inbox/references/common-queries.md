# Common Email Queries

Practical code examples for frequently requested email operations using go-imap v1.

## Prerequisites

All examples assume you have a connected IMAP client:

```go
import (
	"crypto/tls"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func connect(email, password string) (*client.Client, error) {
	c, err := client.DialTLS("imap.gmail.com:993", &tls.Config{})
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

## Get Recent Starred Emails

Fetch the N most recent starred emails from the Starred mailbox.

```go
func GetRecentStarred(c *client.Client, count int) ([]*imap.Message, error) {
	// Select the Starred mailbox
	mbox, err := c.Select("[Gmail]/Starred", false)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return []*imap.Message{}, nil
	}

	// Calculate range to fetch
	fetchCount := uint32(count)
	if mbox.Messages < fetchCount {
		fetchCount = mbox.Messages
	}

	from := mbox.Messages - fetchCount + 1
	to := mbox.Messages

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
		imap.FetchFlags,
	}

	messages := make(chan *imap.Message, count+1)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	// Reverse to get most recent first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// Usage:
// messages, err := GetRecentStarred(c, 5)
// for _, msg := range messages {
//     fmt.Printf("Subject: %s\n", msg.Envelope.Subject)
// }
```

## Get Recent Unread Emails

Fetch the N most recent unread emails from the inbox.

```go
func GetRecentUnread(c *client.Client, count int) ([]*imap.Message, error) {
	// Select inbox
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}

	// Search for unread messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	uids, err := c.UidSearch(criteria)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return []*imap.Message{}, nil
	}

	// Get the most recent N
	fetchCount := count
	if len(uids) < fetchCount {
		fetchCount = len(uids)
	}

	recentUids := uids[len(uids)-fetchCount:]

	// Reverse to get most recent first
	for i, j := 0, len(recentUids)-1; i < j; i, j = i+1, j-1 {
		recentUids[i], recentUids[j] = recentUids[j], recentUids[i]
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(recentUids...)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
	}

	messages := make(chan *imap.Message, count+1)
	done := make(chan error, 1)

	go func() {
		done <- c.UidFetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}
```

## Get Emails by Date Range

Fetch emails within a specific date range.

```go
import "time"

func GetEmailsByDateRange(c *client.Client, start, end time.Time) ([]*imap.Message, error) {
	// Select All Mail for comprehensive search
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Since = start
	criteria.Before = end.Add(24 * time.Hour) // Include the end date

	uids, err := c.UidSearch(criteria)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return []*imap.Message{}, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
	}

	messages := make(chan *imap.Message, len(uids))
	done := make(chan error, 1)

	go func() {
		done <- c.UidFetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}

// Usage:
// start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
// end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
// messages, err := GetEmailsByDateRange(c, start, end)
```

## List Recent Emails with Pagination

Fetch emails in pages for infinite scroll or pagination UI.

```go
type PageResult struct {
	Messages []*imap.Message
	HasMore  bool
	NextPage int
}

func GetEmailsPage(c *client.Client, folder string, page, pageSize int) (*PageResult, error) {
	mbox, err := c.Select(folder, true)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return &PageResult{
			Messages: []*imap.Message{},
			HasMore:  false,
			NextPage: page,
		}, nil
	}

	// Calculate range for this page
	totalMessages := int(mbox.Messages)
	startOffset := totalMessages - (page * pageSize)
	endOffset := startOffset + pageSize - 1

	if startOffset < 0 {
		startOffset = 0
	}

	from := uint32(startOffset + 1)
	to := uint32(endOffset + 1)

	if to > mbox.Messages {
		to = mbox.Messages
	}

	if from > to {
		return &PageResult{
			Messages: []*imap.Message{},
			HasMore:  false,
			NextPage: page,
		}, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
	}

	messages := make(chan *imap.Message, pageSize+1)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	// Reverse to get most recent first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	hasMore := startOffset > 0

	return &PageResult{
		Messages: result,
		HasMore:  hasMore,
		NextPage: page + 1,
	}, nil
}

// Usage:
// result, err := GetEmailsPage(c, "INBOX", 1, 20)  // Page 1, 20 items per page
// if result.HasMore {
//     nextPage, err := GetEmailsPage(c, "INBOX", result.NextPage, 20)
// }
```

## Get Emails from Specific Sender

Fetch recent emails from a specific email address.

```go
func GetEmailsFromSender(c *client.Client, senderEmail string, count int) ([]*imap.Message, error) {
	// Select All Mail
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Header.Add("From", senderEmail)

	uids, err := c.UidSearch(criteria)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return []*imap.Message{}, nil
	}

	// Get the most recent N
	fetchCount := count
	if len(uids) < fetchCount {
		fetchCount = len(uids)
	}

	recentUids := uids[len(uids)-fetchCount:]

	// Reverse to get most recent first
	for i, j := 0, len(recentUids)-1; i < j; i, j = i+1, j-1 {
		recentUids[i], recentUids[j] = recentUids[j], recentUids[i]
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(recentUids...)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
	}

	messages := make(chan *imap.Message, count+1)
	done := make(chan error, 1)

	go func() {
		done <- c.UidFetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}
```

## Search by Subject

Find emails with specific subject keywords.

```go
func SearchBySubject(c *client.Client, keyword string, count int) ([]*imap.Message, error) {
	// Select All Mail
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Header.Add("Subject", keyword)

	uids, err := c.UidSearch(criteria)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return []*imap.Message{}, nil
	}

	// Get the most recent N
	fetchCount := count
	if len(uids) < fetchCount {
		fetchCount = len(uids)
	}

	recentUids := uids[len(uids)-fetchCount:]

	// Reverse to get most recent first
	for i, j := 0, len(recentUids)-1; i < j; i, j = i+1, j-1 {
		recentUids[i], recentUids[j] = recentUids[j], recentUids[i]
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(recentUids...)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
	}

	messages := make(chan *imap.Message, count+1)
	done := make(chan error, 1)

	go func() {
		done <- c.UidFetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}
```

## Get Emails with Attachments

Find emails that have attachments.

```go
func GetEmailsWithAttachments(c *client.Client, count int) ([]*imap.Message, error) {
	// Use Gmail's X-GM-RAW search for "has:attachment"
	// This requires the custom GmailRawSearch from operations.md
	
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}

	// Search for messages with attachments using Gmail search
	searchCmd := &GmailRawSearch{Query: "has:attachment"}
	searchResp := &responses.Search{}
	
	if _, err := c.Execute(searchCmd, searchResp); err != nil {
		return nil, err
	}

	if len(searchResp.Ids) == 0 {
		return []*imap.Message{}, nil
	}

	// Get the most recent N
	fetchCount := count
	if len(searchResp.Ids) < fetchCount {
		fetchCount = len(searchResp.Ids)
	}

	recentUids := searchResp.Ids[len(searchResp.Ids)-fetchCount:]

	// Reverse to get most recent first
	for i, j := 0, len(recentUids)-1; i < j; i, j = i+1, j-1 {
		recentUids[i], recentUids[j] = recentUids[j], recentUids[i]
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(recentUids...)

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchUid,
	}

	messages := make(chan *imap.Message, count+1)
	done := make(chan error, 1)

	go func() {
		done <- c.UidFetch(seqSet, items, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}
```

## Quick Reference: Common Operations

| Operation | Mailbox | Method |
|-----------|---------|--------|
| Get starred emails | `[Gmail]/Starred` | Direct select + fetch |
| Get unread emails | `INBOX` | `UidSearch` with `WithoutFlags: [Seen]` |
| Get all emails | `[Gmail]/All Mail` | Direct select + fetch |
| Get sent emails | `[Gmail]/Sent Mail` | Direct select + fetch |
| Get drafts | `[Gmail]/Drafts` | Direct select + fetch |
| Get trash | `[Gmail]/Trash` | Direct select + fetch |
| Search by date | `[Gmail]/All Mail` | `UidSearch` with `Since`/`Before` |
| Search by sender | `[Gmail]/All Mail` | `UidSearch` with `Header.Add("From", ...)` |
| Search by subject | `[Gmail]/All Mail` | `UidSearch` with `Header.Add("Subject", ...)` |
| Complex search | `[Gmail]/All Mail` | `X-GM-RAW` with Gmail search syntax |
