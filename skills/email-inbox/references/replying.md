# Replying in an Existing Gmail Thread

This reference covers the part that basic SMTP examples usually skip: finding the right Gmail conversation via IMAP, extracting the RFC headers needed for threading, and then sending a reply that stays attached to the existing conversation.

## Key Idea

There are two different threading concepts involved:

- Gmail conversation identity: `X-GM-THRID`
- RFC email threading headers: `Message-ID`, `In-Reply-To`, and `References`

Use both.

- Use IMAP (`X-GM-RAW`, `X-GM-THRID`) to find the right Gmail conversation.
- Use SMTP (`In-Reply-To`, `References`) to build the outgoing reply correctly.

If you only match the subject, threading may still fail.

## Recommended Workflow

1. Select `[Gmail]/All Mail`.
2. Search broadly with `X-GM-RAW`.
3. Fetch candidate envelopes plus `X-GM-THRID`.
4. Pick the intended conversation.
5. Expand the whole conversation with `UID SEARCH X-GM-THRID ...`.
6. Sort by date and choose the latest message.
7. Fetch the latest message's RFC headers.
8. Send SMTP message with `In-Reply-To` and appended `References`.

## IMAP Helpers

### Gmail search commands

```go
type GmailRawSearch struct {
    Query string
}

func (s *GmailRawSearch) Command() *imap.Command {
    return &imap.Command{
        Name:      "UID SEARCH",
        Arguments: []interface{}{imap.RawString(`X-GM-RAW "` + s.Query + `"`)},
    }
}

type GmailThreadSearch struct {
    ThreadID string
}

func (s *GmailThreadSearch) Command() *imap.Command {
    return &imap.Command{
        Name:      "UID SEARCH",
        Arguments: []interface{}{imap.RawString("X-GM-THRID " + s.ThreadID)},
    }
}
```

### Fetch reply-threading headers from a specific UID

This is the go-imap v1 pattern that matters for replies.

```go
import (
    "bufio"
    "net/mail"

    "github.com/emersion/go-imap"
)

type ReplyTarget struct {
    Subject    string
    From       string
    MessageID  string
    References string
}

func FetchReplyTarget(c *client.Client, uid uint32) (*ReplyTarget, error) {
    if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
        return nil, err
    }

    headerSection := &imap.BodySectionName{
        BodyPartName: imap.BodyPartName{
            Specifier: imap.HeaderSpecifier,
            Fields:    []string{"Message-ID", "References", "In-Reply-To"},
        },
        Peek: true,
    }

    seqset := new(imap.SeqSet)
    seqset.AddNum(uid)

    messages := make(chan *imap.Message, 1)
    done := make(chan error, 1)
    go func() {
        done <- c.UidFetch(seqset, []imap.FetchItem{
            imap.FetchEnvelope,
            imap.FetchUid,
            headerSection.FetchItem(),
        }, messages)
    }()

    msg := <-messages
    if err := <-done; err != nil {
        return nil, err
    }
    if msg == nil {
        return nil, fmt.Errorf("uid %d not found", uid)
    }

    target := &ReplyTarget{}
    if msg.Envelope != nil {
        target.Subject = msg.Envelope.Subject
        if len(msg.Envelope.From) > 0 {
            target.From = msg.Envelope.From[0].Address()
        }
    }

    r := msg.GetBody(headerSection)
    if r == nil {
        return nil, fmt.Errorf("missing header body for uid %d", uid)
    }

    parsed, err := mail.ReadMessage(bufio.NewReader(r))
    if err != nil {
        return nil, err
    }

    target.MessageID = strings.TrimSpace(parsed.Header.Get("Message-ID"))
    target.References = strings.TrimSpace(parsed.Header.Get("References"))

    return target, nil
}
```

Why this exact shape matters:
- In go-imap v1, the requested header fields live under `BodyPartName.Fields`.
- Set `Specifier: imap.HeaderSpecifier`.
- Then request `headerSection.FetchItem()`.
- Parse the returned header block with `net/mail`.

## Expanding a Gmail Conversation First

Once you know a Gmail thread ID, fetch all messages in the conversation, then choose the newest message as the reply anchor.

```go
const FetchGmailThreadID imap.FetchItem = "X-GM-THRID"

type ThreadMessage struct {
    UID      uint32
    Date     time.Time
    Subject  string
    ThreadID string
}

func GetThreadMessages(c *client.Client, threadID string) ([]ThreadMessage, error) {
    if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
        return nil, err
    }

    searchResp := &responses.Search{}
    if _, err := c.Execute(&GmailThreadSearch{ThreadID: threadID}, searchResp); err != nil {
        return nil, err
    }
    if len(searchResp.Ids) == 0 {
        return nil, fmt.Errorf("thread %s not found", threadID)
    }

    seqset := new(imap.SeqSet)
    for _, uid := range searchResp.Ids {
        seqset.AddNum(uid)
    }

    messages := make(chan *imap.Message, len(searchResp.Ids)+1)
    done := make(chan error, 1)
    go func() {
        done <- c.UidFetch(seqset, []imap.FetchItem{
            imap.FetchEnvelope,
            imap.FetchUid,
            imap.FetchInternalDate,
            FetchGmailThreadID,
        }, messages)
    }()

    var out []ThreadMessage
    for msg := range messages {
        tm := ThreadMessage{UID: msg.Uid, Date: msg.InternalDate}
        if msg.Envelope != nil {
            tm.Subject = msg.Envelope.Subject
            tm.Date = msg.Envelope.Date
        }
        tm.ThreadID = fmt.Sprintf("%v", msg.Items[FetchGmailThreadID])
        out = append(out, tm)
    }

    if err := <-done; err != nil {
        return nil, err
    }

    sort.Slice(out, func(i, j int) bool {
        return out[i].Date.Before(out[j].Date)
    })
    return out, nil
}
```

## Building the Outgoing Reply

When replying, append the current message ID to the existing references chain.

```go
func BuildReplyMessage(from, to, subject, body string, target *ReplyTarget) string {
    msgID := fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), from)

    if !strings.HasPrefix(strings.ToLower(subject), "re:") {
        subject = "Re: " + subject
    }

    refs := strings.TrimSpace(strings.TrimSpace(target.References + " " + target.MessageID))
    refs = strings.Join(strings.Fields(refs), " ")

    headers := []string{
        fmt.Sprintf("From: %s", from),
        fmt.Sprintf("To: %s", to),
        fmt.Sprintf("Subject: %s", subject),
        fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
        fmt.Sprintf("Message-ID: %s", msgID),
        fmt.Sprintf("In-Reply-To: %s", target.MessageID),
        fmt.Sprintf("References: %s", refs),
        "MIME-Version: 1.0",
        "Content-Type: text/plain; charset=UTF-8",
    }

    return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
```

Why append `References` instead of replacing it:
- `In-Reply-To` points at the immediate parent.
- `References` preserves the full RFC thread chain.
- Reusing only the parent `Message-ID` is often enough for simple cases, but appending the full chain is more robust.

## End-to-End Pattern

```go
func ReplyInExistingThread(c *client.Client, smtpClient *smtp.Client, from, to, rawQuery, body string) error {
    // 1. Find candidate messages with a broad Gmail search.
    if _, err := c.Select("[Gmail]/All Mail", true); err != nil {
        return err
    }

    searchResp := &responses.Search{}
    if _, err := c.Execute(&GmailRawSearch{Query: rawQuery}, searchResp); err != nil {
        return err
    }
    if len(searchResp.Ids) == 0 {
        return fmt.Errorf("no messages matched %q", rawQuery)
    }

    // 2. Pick one candidate UID, fetch its reply headers, and build message.
    // In production, prefer fetching X-GM-THRID first and expanding to the whole thread.
    latestUID := searchResp.Ids[len(searchResp.Ids)-1]

    target, err := FetchReplyTarget(c, latestUID)
    if err != nil {
        return err
    }

    msg := BuildReplyMessage(from, to, target.Subject, body, target)

    if err := smtpClient.Mail(from, nil); err != nil {
        return err
    }
    if err := smtpClient.Rcpt(to, nil); err != nil {
        return err
    }
    w, err := smtpClient.Data()
    if err != nil {
        return err
    }
    if _, err := w.Write([]byte(msg)); err != nil {
        return err
    }
    return w.Close()
}
```

## Search Robustness Tips

For support threads and order conversations, do not rely on one ultra-precise `X-GM-RAW` query.

Prefer this sequence:
- broad sender/date query
- fetch envelopes and thread IDs
- identify the correct conversation locally
- expand by `X-GM-THRID`
- reply to the latest message in that conversation

This avoids false negatives from Gmail IMAP search parsing and avoids false positives from ambiguous terms.

## Common Failure Modes

### Narrow Gmail query returns zero unexpectedly

Reason:
- Gmail IMAP `X-GM-RAW` is powerful, but exact field/quote combinations can be finicky.

Fix:
- broaden the query
- add sender/date restrictions
- filter locally after fetch

### Reply sends, but does not stay in thread

Reason:
- outgoing message used only the same subject, without `In-Reply-To` / `References`

Fix:
- fetch the original message headers from IMAP
- set both `In-Reply-To` and `References`

### Wrong message chosen as reply anchor

Reason:
- candidate search found a related message, survey email, auto-response, or older branch

Fix:
- expand by `X-GM-THRID`
- sort messages by date
- choose the newest relevant message in that Gmail conversation
