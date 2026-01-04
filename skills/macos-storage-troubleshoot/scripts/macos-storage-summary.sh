#!/bin/bash
# macos-storage-summary.sh
# Quick summary of macOS storage usage

echo "=== DISK OVERVIEW ==="
df -h /

echo ""
echo "=== APFS CONTAINER ==="
diskutil apfs list | grep -E "(Capacity|Name:.*Case|Volume.*Role)"

echo ""
echo "=== TOP DIRECTORIES IN HOME ==="
du -h -d 1 ~ 2>/dev/null | sort -hr | head -15

echo ""
echo "=== MESSAGES ATTACHMENTS ==="
du -sh ~/Library/Messages/Attachments 2>/dev/null || echo "Cannot access (need Full Disk Access?)"

echo ""
echo "=== LARGE FILES IN MESSAGES (>100MB) ==="
find ~/Library/Messages/Attachments -type f -size +100M -exec ls -lh {} \; 2>/dev/null | wc -l | xargs echo "Count:"
