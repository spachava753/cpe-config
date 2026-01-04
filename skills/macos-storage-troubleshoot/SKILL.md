---
name: macos-storage-troubleshoot
description: Troubleshoot macOS storage issues, especially "System Data" that appears large in System Settings. Use when users ask about disk space, storage usage, what's taking up space, System Data, or need to find and clean up large files. Covers APFS volumes, Messages attachments, Cryptexes, caches, and hidden system directories.
---

# macOS Storage Troubleshooting

## Quick Start

1. Check APFS volume usage vs `du` output - discrepancy reveals "hidden" data
2. Grant Full Disk Access to Terminal for complete visibility
3. Check Messages attachments (common culprit for large "System Data")

## Key Insight: The Discrepancy Method

macOS "System Data" often appears larger than expected. Find what's hidden:

```bash
# APFS reports this much for Data volume:
diskutil info disk2s5 | grep "Volume Used Space"

# But du can only see this much:
sudo du -sh /System/Volumes/Data

# DISCREPANCY = Hidden "System Data"
```

## Common Storage Locations

### User-Accessible (no sudo)
| Location | Command |
|----------|---------|
| Home directory | `du -sh ~/*` |
| Library caches | `du -sh ~/Library/Caches` |
| Application Support | `du -sh ~/Library/Application\ Support` |
| Downloads | `du -sh ~/Downloads` |

### Requires Full Disk Access
| Location | What it contains |
|----------|------------------|
| `/System/Volumes/Data/.Spotlight-V100` | Spotlight search index |
| `/System/Volumes/Data/.DocumentRevisions-V100` | File version history |
| `/System/Volumes/Data/.fseventsd` | Filesystem events |
| `~/Library/Messages/Attachments` | iMessage attachments |

### System Volumes (not reducible)
| Volume | Typical Size | Notes |
|--------|--------------|-------|
| Preboot/Cryptexes | 20-30 GB | Security updates, Safari, AI models |
| Sealed System | 12-15 GB | macOS itself (immutable) |
| Recovery | 1-2 GB | Recovery partition |

## Messages Attachments Deep Dive

Messages attachments are a common culprit. They often get categorized as "System Data" in Storage settings.

### Check Messages size
```bash
sudo du -sh ~/Library/Messages/Attachments
```

### Find large files
```bash
find ~/Library/Messages/Attachments -type f -size +100M -exec ls -lh {} \; | sort -k5 -hr | head -20
```

### Query database for attachment details
See [references/messages-queries.md](references/messages-queries.md) for SQLite queries to:
- Find largest attachments with conversation info
- Look up contacts for phone numbers
- Identify orphaned attachments
- Handle integer overflow for files >2GB

### Safe deletion workflow
1. Delete attachments via Messages app UI (safest)
2. Or delete files directly: `rm ~/Library/Messages/Attachments/path/to/file`
3. For iCloud sync of deletions, see [references/messages-queries.md](references/messages-queries.md)

## Cleanup Commands

### Safe cleanups
```bash
# Clear Homebrew cache
brew cleanup --prune=all

# Clear Go caches
go clean -cache
go clean -modcache

# Clear npm cache
npm cache clean --force

# Rebuild Spotlight index (clears .Spotlight-V100)
sudo mdutil -E /

# Delete old iOS backups
# System Settings → General → Storage → iOS Files
```

### Time Machine snapshots
```bash
# List local snapshots
tmutil listlocalsnapshots /

# Delete all local snapshots
sudo tmutil deletelocalsnapshots /
```

## Full Disk Access Setup

Required to access protected directories:
1. System Settings → Privacy & Security → Full Disk Access
2. Add Terminal (or iTerm, etc.)
3. Restart terminal

## Troubleshooting Checklist

1. [ ] Check overall disk usage: `df -h /`
2. [ ] Check APFS container: `diskutil apfs list`
3. [ ] Compare APFS vs du for Data volume
4. [ ] Grant Full Disk Access if needed
5. [ ] Check Messages attachments size
6. [ ] Check ~/Library breakdown: `sudo du -h -d 1 ~/Library | sort -hr | head -20`
7. [ ] Check for Time Machine snapshots
8. [ ] Review Cryptexes/Preboot (usually not reducible)
