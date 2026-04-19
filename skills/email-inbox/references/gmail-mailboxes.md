# Gmail Mailbox Structure

Gmail uses a special IMAP folder structure. Understanding these mailboxes is essential for effective email operations.

Gmail also returns special-use mailbox attributes in normal `LIST` responses. Prefer those standard special-use attributes over legacy `XLIST`, which Gmail deprecated.

## Standard Gmail Mailboxes

| Mailbox Path | Description | Use Case |
|--------------|-------------|----------|
| `INBOX` | Main inbox | Reading new mail, checking unread counts |
| `[Gmail]/All Mail` | Contains all emails | Comprehensive searches, thread operations |
| `[Gmail]/Starred` | Starred/flagged emails | Getting starred items |
| `[Gmail]/Important` | Emails marked as important | Priority email handling |
| `[Gmail]/Sent Mail` | Sent emails | Outbox/sent history |
| `[Gmail]/Drafts` | Draft emails | Working with drafts |
| `[Gmail]/Trash` | Deleted emails | Trash management |
| `[Gmail]/Spam` | Spam/junk emails | Spam handling |

## When to Use Each Mailbox

### INBOX
Use for:
- Reading new emails
- Checking unread counts
- Basic email operations on recent mail
- Operations that should only affect inbox items

```go
mbox, err := c.Select("INBOX", false)
```

### [Gmail]/All Mail
Use for:
- Comprehensive searches across all emails
- Thread operations (threads may span multiple folders)
- Finding emails regardless of labels
- Global date range searches

```go
mbox, err := c.Select("[Gmail]/All Mail", false)
```

**Note:** This mailbox contains ALL emails, including sent, drafts, and spam. Use specific search criteria to filter.

### [Gmail]/Starred
Use for:
- Getting starred/flagged emails
- Managing starred items
- Quick access to important emails

```go
mbox, err := c.Select("[Gmail]/Starred", false)
```

**Tip:** This is the most reliable way to get starred emails. Using `X-GM-RAW` search with `is:starred` may not work as expected in all cases.

### [Gmail]/Sent Mail
Use for:
- Accessing sent emails
- Checking sent history
- Thread operations involving sent messages

```go
mbox, err := c.Select("[Gmail]/Sent Mail", false)
```

### [Gmail]/Drafts
Use for:
- Working with draft emails
- Saving drafts
- Retrieving unsent messages

```go
mbox, err := c.Select("[Gmail]/Drafts", false)
```

### [Gmail]/Trash
Use for:
- Permanently deleting emails
- Restoring deleted emails
- Managing trash contents

```go
mbox, err := c.Select("[Gmail]/Trash", false)
```

## User-Created Labels

Gmail labels appear as mailboxes. Labels with spaces are wrapped in quotes.

```go
// Label without spaces
mbox, err := c.Select("Work", false)

// Label with spaces (quoted)
mbox, err := c.Select(`"Work Projects"`, false)
```

## Important Considerations

### 1. Read-Only vs Read-Write Mode

```go
// Read-only (faster, no lock)
mbox, err := c.Select("INBOX", true)

// Read-write (required for modifications)
mbox, err := c.Select("INBOX", false)
```

Use read-only (`true`) for:
- Listing emails
- Searching
- Reading content

Use read-write (`false`) for:
- Adding/removing labels
- Marking as read/unread
- Deleting emails
- Moving between folders

### 2. Message Counts

```go
mbox, err := c.Select("[Gmail]/All Mail", true)
if err != nil {
    return err
}

fmt.Printf("Total messages: %d\n", mbox.Messages)
fmt.Printf("Recent messages: %d\n", mbox.Recent)
fmt.Printf("Unseen messages: %d\n", mbox.Unseen)
```

### 3. Cross-Mailbox Threads

Gmail threads can span multiple mailboxes (e.g., an email in INBOX with replies in Sent). When working with threads:

- Always use `[Gmail]/All Mail` for thread searches
- Thread IDs are consistent across mailboxes
- Messages in a thread may have different labels

```go
// Good: Search for thread in All Mail
if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
    return err
}
// ... search by thread ID

// Bad: Thread might be split across mailboxes
if _, err := c.Select("INBOX", true); err != nil {
    return err
}
// ... might miss messages in Sent
```

### 4. Special Mailbox Behavior

| Mailbox | Behavior |
|---------|----------|
| `[Gmail]/All Mail` | Contains everything; messages appear here even when in other folders |
| `[Gmail]/Starred` | Messages here are also in All Mail; starring adds the `\\Flagged` flag |
| `[Gmail]/Important` | Gmail's importance algorithm; can be manually adjusted |
| `[Gmail]/Trash` | Messages here are excluded from All Mail searches by default |

## Common Patterns

### Pattern 1: Get Recent Emails from Specific Folder

```go
func GetRecentFromFolder(c *client.Client, folder string, count int) ([]*imap.Message, error) {
	mbox, err := c.Select(folder, true)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return []*imap.Message{}, nil
	}

	fetchCount := uint32(count)
	if mbox.Messages < fetchCount {
		fetchCount = mbox.Messages
	}

	from := mbox.Messages - fetchCount + 1
	to := mbox.Messages

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// ... fetch messages
}
```

### Pattern 2: Search Across All Mail

```go
func SearchAllMail(c *client.Client, criteria *imap.SearchCriteria) ([]uint32, error) {
	if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
		return nil, err
	}
	return c.UidSearch(criteria)
}
```

### Pattern 3: Work with Starred Items

```go
func GetStarredCount(c *client.Client) (uint32, error) {
	mbox, err := c.Select("[Gmail]/Starred", true)
	if err != nil {
		return 0, err
	}
	return mbox.Messages, nil
}
```
