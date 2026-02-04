---
name: email-inbox
description: Manage user's email inbox via Gmail IMAP and send emails via SMTP. Use when the user asks to check email, search emails, read messages, organize inbox, manage labels, delete emails, work with email threads/conversations, or send/compose/reply to emails. Requires GMAIL_APP_PASSWORD environment variable.
---

# Email Management

Manage emails via Gmail IMAP (`github.com/emersion/go-imap` v1) and send emails via SMTP (`github.com/emersion/go-smtp`).

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
- **Send an email** (new, reply, with attachments)
- **Compose or draft an email**
- **Reply to an email thread**

## Prerequisites

**Gmail App Password**: Requires `GMAIL_APP_PASSWORD` environment variable.
- User must enable 2FA on Google account
- Generate app password at https://myaccount.google.com/apppasswords

**Gmail Address**: Requires `GMAIL_ADDRESS` environment variable for sending.

## Reference Documentation

- **[Gmail Mailbox Structure](./references/gmail-mailboxes.md)** - Understanding Gmail's IMAP folder structure and when to use each mailbox
- **[Common Email Queries](./references/common-queries.md)** - Practical examples for frequently requested operations (starred, unread, date ranges, pagination, etc.)
- **[Thread Operations](./references/operations.md)** - Complete code examples for Gmail thread management
- **[Sending Email](./references/sending.md)** - Complete code examples for sending emails via SMTP

## Operations

### Reading Email (IMAP)
See `references/operations.md` for complete code examples.

- **List Threads** - Recent conversations with subjects, senders, dates
- **Get Thread** - All messages in a conversation by thread ID
- **Search** - Find emails using Gmail search syntax

### Organizing Email (IMAP)
- **Add Label** - Apply label to a thread
- **Remove Label** - Remove label from a thread

### Deleting Email (IMAP)
- **Delete Thread** - Move thread to Trash

### Sending Email (SMTP)
See `references/sending.md` for complete code examples.

- **Send Text Email** - Plain text email
- **Send HTML Email** - Email with HTML content
- **Send to Multiple** - To, CC, BCC recipients
- **Send with Attachment** - Email with file attachments
- **Reply to Email** - Reply with In-Reply-To header

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

## Gmail SMTP Settings

| Setting | Value |
|---------|-------|
| Server | `smtp.gmail.com` |
| Port | `465` (SSL/TLS) |
| Auth | App Password (not regular password) |

**Sending Limits:**
- Free Gmail: 500 emails/day
- Google Workspace: 2,000 emails/day

## Technical Notes

**IMAP (Reading):**
- Uses `github.com/emersion/go-imap` v1.2.1 (not v2)
- Gmail extensions: X-GM-THRID (thread ID), X-GM-LABELS (labels), X-GM-RAW (search)
- Select `[Gmail]/All Mail` for comprehensive operations
- Thread ID is a 64-bit integer unique to each conversation

**SMTP (Sending):**
- Uses `github.com/emersion/go-smtp` + `github.com/emersion/go-sasl`
- Uses TLS on port 465 (direct SSL connection)
- Build RFC 5322 compliant messages manually
- Use `In-Reply-To` and `References` headers for replies
