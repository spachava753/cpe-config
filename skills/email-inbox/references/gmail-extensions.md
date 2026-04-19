# Gmail IMAP Extensions and Behavior

Gmail supports normal IMAP operations from RFC 3501. This reference covers the Gmail-specific extensions and altered behavior that matter when building tooling on top of Gmail IMAP.

Primary source: https://developers.google.com/workspace/gmail/imap/imap-extensions

## What Gmail Adds

- `X-GM-EXT-1` capability advertises Gmail's extension set.
- Special-use mailbox attributes are returned by normal `LIST` responses.
- `X-GM-RAW` exposes Gmail web-style search syntax through IMAP.
- `X-GM-MSGID` exposes Gmail's stable per-message ID.
- `X-GM-THRID` exposes Gmail's stable thread ID.
- `X-GM-LABELS` exposes Gmail labels and lets clients modify them.

These are additions to standard IMAP, not replacements for it.

## Capability and Client ID

Gmail advertises extension support via `CAPABILITY`. Look for `X-GM-EXT-1`.

Google also recommends sending an IMAP `ID` command with client name/version/contact information. This is helpful, but not required for ordinary mailbox operations.

Practical guidance for this skill:
- Do not assume Gmail-specific features exist on non-Gmail IMAP servers.
- For Gmail-specific workflows, select Gmail first and then use `X-GM-*` commands.
- Keep standard IMAP fallbacks in mind when Gmail-only features are unnecessary.

## Special-Use Mailboxes

Gmail supports RFC 6154 special-use mailboxes. In practice, `LIST` responses identify folders like:
- `INBOX`
- `[Gmail]/All Mail`
- `[Gmail]/Drafts`
- `[Gmail]/Sent Mail`
- `[Gmail]/Spam`
- `[Gmail]/Trash`
- `[Gmail]/Starred`
- `[Gmail]/Important`

### Important notes

- Prefer special-use `LIST` discovery over legacy `XLIST`.
- `XLIST` is deprecated.
- `[Gmail]/All Mail` is the safest mailbox for global searches and thread work.
- Threads can span inbox, sent, archived, and labeled views.

## X-GM-RAW Search

`X-GM-RAW` lets you use Gmail search syntax in IMAP:

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
```

Examples:
- `from:alice@example.com`
- `subject:invoice after:2026/01/01`
- `has:attachment in:unread`
- `label:Work`

### Practical caveats

`X-GM-RAW` is powerful, but narrow compound queries can be brittle over IMAP. When a precise query unexpectedly returns zero:
- retry with a broader sender/date query
- fetch a small candidate set
- filter subjects, body snippets, or thread IDs locally

This is especially useful for support threads, order lookups, and ambiguous terms like `daylight`.

## X-GM-MSGID

`X-GM-MSGID` is Gmail's stable message identifier across folders.

Use it when you need to:
- identify a single message reliably
- deduplicate the same message seen through multiple labels/mailboxes
- jump back to a specific Gmail message from a stored ID

```go
const FetchGmailMessageID imap.FetchItem = "X-GM-MSGID"

items := []imap.FetchItem{imap.FetchEnvelope, FetchGmailMessageID}
```

Search by Gmail message ID:

```go
type GmailMsgIDSearch struct {
    MessageID string
}

func (s *GmailMsgIDSearch) Command() *imap.Command {
    return &imap.Command{
        Name:      "UID SEARCH",
        Arguments: []interface{}{imap.RawString("X-GM-MSGID " + s.MessageID)},
    }
}
```

## X-GM-THRID

`X-GM-THRID` is Gmail's stable thread identifier.

Use it when you need to:
- group messages the way Gmail's UI groups them
- fetch all messages in a conversation across inbox/sent/archive
- apply labels or deletes to every message in a Gmail thread

```go
const FetchGmailThreadID imap.FetchItem = "X-GM-THRID"

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

Important distinction:
- `X-GM-THRID` is Gmail conversation identity.
- `Message-ID` / `In-Reply-To` / `References` are RFC email-threading headers.

For a reply that lands in the existing Gmail conversation reliably, use both:
- use IMAP and `X-GM-THRID` or `X-GM-RAW` to find the target conversation
- use SMTP with `In-Reply-To` and `References` to build the outgoing reply correctly

## X-GM-LABELS

Gmail labels are exposed via `X-GM-LABELS`.

```go
const FetchGmailLabels imap.FetchItem = "X-GM-LABELS"
```

You can fetch labels:

```go
items := []imap.FetchItem{imap.FetchEnvelope, FetchGmailLabels}
```

You can also add or remove labels with `UID STORE` plus `+X-GM-LABELS` or `-X-GM-LABELS`.

Practical guidance:
- system labels live under `[Gmail]` and should be treated as reserved
- custom labels may need quoting when they contain spaces
- a message can appear in multiple label views without being duplicated logically

## Recommended Workflow Patterns

### Read/search workflows

- Use standard IMAP mailbox selection.
- Prefer `[Gmail]/All Mail` for cross-thread or cross-label searches.
- Use `X-GM-RAW` for Gmail-style search syntax.
- Fetch `X-GM-MSGID` or `X-GM-THRID` when you need stable identity.

### Thread workflows

- Search candidate messages with `X-GM-RAW`.
- Fetch `X-GM-THRID` for candidates.
- Search by `X-GM-THRID` to expand to the whole conversation.
- Sort by date before choosing the latest message.

### Reply workflows

- Find the target message or thread via IMAP.
- Fetch the latest message's RFC headers (`Message-ID`, `References`, optionally `In-Reply-To`).
- Build the outgoing email with `In-Reply-To` and an appended `References` chain.
- Send through SMTP.

See `replying.md` for a complete thread-aware reply pattern.
