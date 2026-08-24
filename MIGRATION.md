# komp P1 Migration Checklist

## 1. Install binary

```bash
cd /Users/rd/Scripts/Files/compression/komp
go build -o /opt/homebrew/bin/komp .
komp --version   # should print nothing (no version flag yet) or help
```

## 2. Rewire 3 KM palette entries

Replace each codec macro's action with **Execute Shell Script**:

```bash
komp --finder -f zstd --delete
komp --finder -f zip --delete
komp --finder -f tar
```

Keep original hotkeys. Test: select files in Finder → hit hotkey → verify archive + originals behavior.

## 3. Smoke-test round-trip

- `komp add --last` after a create → adds file → notification on large job
- `komp ls <archive>` → contents
- `komp -d 7z ~/tmp/test` → explicit path, delete originals

## 4. On confirmation, retire old layer

```bash
cd /Users/rd/Scripts/Files/compression
mkdir -p _attic
mv keep remove shell archive-selection.sh archive-selection.sh.bak \
   compress-and-remove.sh _attic/ 2>/dev/null
# macros/*.kmmacros stay until P3 retires convert/img duties
```

## 5. Update Obsidian vault

- Create `Files/Compression/Komp Toolkit.md` (Smart Sentence Case)
- Mirror from repo README + any vault-specific notes
- Ensure `_docs` symlink points to vault folder
