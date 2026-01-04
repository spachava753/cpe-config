# Messages Database Queries

## Database Location
```
~/Library/Messages/chat.db
```

## Find Largest Attachments with Conversation Info

```sql
SELECT 
    attachment.filename,
    ROUND(attachment.total_bytes / 1024.0 / 1024.0 / 1024.0, 2) as size_gb,
    COALESCE(chat.chat_identifier, 'Unknown') as conversation,
    COALESCE(chat.display_name, '') as display_name,
    datetime(message.date/1000000000 + 978307200, 'unixepoch', 'localtime') as date_sent,
    CASE WHEN message.is_from_me = 1 THEN 'Sent' ELSE 'Received' END as direction
FROM attachment
LEFT JOIN message_attachment_join ON attachment.ROWID = message_attachment_join.attachment_id
LEFT JOIN message ON message_attachment_join.message_id = message.ROWID
LEFT JOIN chat_message_join ON message.ROWID = chat_message_join.message_id
LEFT JOIN chat ON chat_message_join.chat_id = chat.ROWID
WHERE attachment.total_bytes > 100000000
ORDER BY attachment.total_bytes DESC
LIMIT 30;
```

## Storage by Conversation

```sql
SELECT 
    COALESCE(chat.chat_identifier, 'Unknown') as contact,
    COALESCE(chat.display_name, '') as name,
    COUNT(attachment.ROWID) as num_attachments,
    ROUND(SUM(attachment.total_bytes) / 1024.0 / 1024.0 / 1024.0, 2) as total_gb
FROM attachment
LEFT JOIN message_attachment_join ON attachment.ROWID = message_attachment_join.attachment_id
LEFT JOIN message ON message_attachment_join.message_id = message.ROWID
LEFT JOIN chat_message_join ON message.ROWID = chat_message_join.message_id
LEFT JOIN chat ON chat_message_join.chat_id = chat.ROWID
GROUP BY chat.chat_identifier
ORDER BY SUM(attachment.total_bytes) DESC
LIMIT 15;
```

## Integer Overflow Issue

Files >2GB have negative `total_bytes` due to 32-bit signed integer overflow.

### Find overflow entries
```sql
SELECT ROWID, filename, total_bytes 
FROM attachment 
WHERE total_bytes < 0;
```

### Fix overflow display (set to 0)
```sql
UPDATE attachment SET total_bytes = 0 WHERE total_bytes < 0;
```

## Look Up Contacts for Phone Numbers

Contacts are in AddressBook databases:
```bash
find ~/Library/Application\\ Support/AddressBook -name "*.abcddb"
```

### Query contacts
```sql
-- Run against each AddressBook-v22.abcddb file
SELECT 
    COALESCE(r.ZFIRSTNAME, '') || ' ' || COALESCE(r.ZLASTNAME, '') as name,
    p.ZFULLNUMBER as phone
FROM ZABCDRECORD r
JOIN ZABCDPHONENUMBER p ON r.Z_PK = p.ZOWNER
WHERE p.ZFULLNUMBER LIKE '%NUMBERHERE%';
```

## Syncing Deletions to iCloud

After manually deleting attachment files, the database entries remain. To sync deletions to iCloud:

### Check for orphaned entries (files deleted but DB entry remains)
```sql
SELECT ROWID, filename, guid, ck_record_id
FROM attachment
WHERE total_bytes = 0 OR total_bytes < 0;
```

### Queue deletions for iCloud sync
```sql
-- This mimics what the delete trigger does
INSERT INTO sync_deleted_attachments (guid, recordID)
SELECT guid, ck_record_id 
FROM attachment 
WHERE ROWID IN (/* list of ROWIDs to delete */);
```

**Warning:** Do NOT use DELETE on the attachment table directly - it has triggers that require Messages app functions. Use UPDATE to fix values, or INSERT into sync queue.

## Messages Database Triggers

The attachment table has these triggers:
- `add_to_sync_deleted_attachments` - Queues deletion for iCloud sync
- `before_delete_on_attachment` - Pre-deletion cleanup
- `after_delete_on_attachment` - Deletes file from disk

These triggers call custom functions only available when Messages.app is running, so external DELETE queries will fail.
