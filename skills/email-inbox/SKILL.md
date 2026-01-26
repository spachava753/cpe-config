---
name: email-inbox
description: Manage user's email inbox via Gmail IMAP. Use when the user asks to check email, search emails, read messages, organize inbox, manage labels, delete emails, or work with email threads/conversations. Requires GMAIL_APP_PASSWORD environment variable.
---

# Email Inbox Management

Manage emails and inbox via Gmail IMAP using `github.com/emersion/go-imap` v1.

## When to Use This Skill

Use when the user asks to:
- Check or read their email/inbox
- Search for emails (by sender, subject, date, unread, etc.)
- List recent emails or conversations
- Read a specific email thread
- Organize emails with labels
- Delete emails or threads
- Find emails with attachments
- Check unread count

## Prerequisites

**Gmail App Password**: Requires `GMAIL_APP_PASSWORD` environment variable.
- User must enable 2FA on Google account
- Generate app password at https://myaccount.google.com/apppasswords

## Operations

See `references/operations.md` for complete code examples.

### Reading Email
- **List Threads** - Recent conversations with subjects, senders, dates
- **Get Thread** - All messages in a conversation by thread ID
- **Search** - Find emails using Gmail search syntax

### Organizing Email  
- **Add Label** - Apply label to a thread
- **Remove Label** - Remove label from a thread

### Deleting Email
- **Delete Thread** - Move thread to Trash

## Gmail Search Syntax (X-GM-RAW)

| Query | Description |
|-------|-------------|
| `is:unread` | Unread emails |
| `is:starred` | Starred emails |
| `from:name@example.com` | From specific sender |
| `to:name@example.com` | To specific recipient |
| `subject:keyword` | Subject contains word |
| `has:attachment` | Has attachments |
| `larger:5M` | Larger than 5MB |
| `after:2024/01/01` | After date |
| `before:2024/12/31` | Before date |
| `label:Work` | Has label |
| `in:inbox` | In inbox |

Combine queries: `is:unread from:github has:attachment`

## Technical Notes

- Uses `github.com/emersion/go-imap` v1.2.1 (not v2)
- Gmail extensions: X-GM-THRID (thread ID), X-GM-LABELS (labels), X-GM-RAW (search)
- Select `[Gmail]/All Mail` for comprehensive operations
- Thread ID is a 64-bit integer unique to each conversation
