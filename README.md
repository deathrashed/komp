# komp

Unified archive & image toolkit for macOS. Compress, add, list, extract, clean, test, and convert archives; build and extract dmg/iso/pkg images. One binary, 17 codecs, interactive TTY picker, and native Finder selection.

## Install

```sh
go install github.com/<owner>/komp@latest
```

> **Module path swap:** The repo currently uses `module komp` for local development. Before publishing, update `go.mod` to `module github.com/<owner>/komp` and push. The install line above assumes the published module path.

## Quickstart

```sh
# Create a zip from files (non-TTY / piped)
komp -f zip notes.txt screenshot.png

# Add files to an existing archive
komp add project.zip newfile.go

# List archive contents
komp ls project.zip

# Extract an archive
komp un project.zip

# Strip junk members from an archive
komp clean --groups macos,vcs project.zip

# Test archive integrity
komp t project.zip

# Interactive mode — omit files/format in a terminal to get a TUI picker
komp
```

## Commands

| Command | Description |
|---------|-------------|
| `komp [format] [files...]` | Compress files/folders (default command) |
| `komp add <archive> [files...]` | Add files to an existing archive |
| `komp ls <archive>` | List archive contents |
| `komp un <archive>...` | Extract archives (and disk images) |
| `komp clean <archive>...` | Strip junk members from archives |
| `komp t <archive>...` | Test archive integrity |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | — | codec name (also accepted as first bare arg) |
| `--delete` | `-d` | false | remove originals after success |
| `--output` | `-o` | input dir | output directory |
| `--level` | `-L` | `normal` | `fast` \| `normal` \| `max` |
| `--separate` | | false | one archive per input |
| `--each` | | false | streams: compress each input separately |
| `--finder` | | false | use current Finder selection as input |
| `--dry-run` | | false | print plan, touch nothing |
| `--backup` | | false | back up overwritten targets |

### un flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | — | destination directory |
| `--here` | | false | extract into current directory |
| `--overwrite` | | false | replace existing files |

### clean flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--groups` | | — | comma list: `macos,windows,vcs,hidden` |
| `--dry-run` | | false | show what would be removed |
| `--backup` | | false | save `<archive>.bak` first |

## Junk groups

`komp clean` strips known junk from archives. Select groups with `--groups`:

| Group | Matches |
|-------|---------|
| `macos` | `.DS_Store`, `._*`, `__MACOSX/`, `.Spotlight-V100`, `.Trashes`, `.fseventsd`, `.LSOverride`, `Icon\r` |
| `windows` | `Thumbs.db`, `desktop.ini`, `$RECYCLE.BIN/` |
| `vcs` | `.git/`, `.svn/`, `.hg/` (any depth) |
| `hidden` | dotfiles not already claimed by another group |

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | success |
| 1 | runtime / environment error (no input, aborted, osascript failure) |
| 2 | usage error (unknown format, missing binary, no format in non-interactive mode) |
| 3 | execution error (compressor / adder / lister / cleaner failed) |
| 4 | integrity failure (verifier detected corruption) |

## Codecs

17 codecs out of the box. Install missing ones via Homebrew:

| Codec | Brew formula |
|-------|-------------|
| zip, tar, gzip, bzip2 | *(system)* |
| 7z | `sevenzip` |
| rar | `rar` |
| zpaq | `zpaq` |
| aar | `aar` |
| xz, lzma | `xz` |
| lzip | `lzip` |
| zstd | `zstd` |
| lz4 | `lz4` |
| brotli | `brotli` |
| lzo | `lzop` |
| snappy | `snzip` |
| lrzip | `lrzip` |

## KM adapter recipes

Keyboard Maestro macros can call `komp` directly. Because `komp --finder` reads the active Finder selection via osascript, a macro only needs to trigger the right flags.

| Macro name | Command | Notes |
|------------|---------|-------|
| Compress selection → zstd | `komp --finder -f zstd --delete` | Deletes originals after success |
| Compress selection → zip | `komp --finder -f zip` | Keeps originals |
| Add to most recent archive | `komp add --last --finder` | Requires at least one prior archive |
| Choose recent archive interactively | `komp add --recent --finder` | TTY-only picker |
| List recent archive | `komp ls --recent` | TTY-only picker |
| Extract archive | `komp un --finder` | Prompts for destination |
| Clean junk from archive | `komp clean --groups macos,vcs <archive>` | Explicit archive required |
| Test archive integrity | `komp t <archive>` | Reports OK or corruption |
| Batch compress (one per file) | `komp --finder -f zip --separate --output ~/Archives` | Creates `~/Archives/<name>.zip` per selection |
| Dry-run plan | `komp --finder -f 7z --dry-run` | Prints plan to stderr, touches nothing |

## Apple-shadowing

`komp` collides with nothing on macOS. No system tool, built-in command, or common framework uses the name `komp`. Safe to drop into `$PATH` without shadowing.

## Roadmap

| Phase | Status | Features |
|-------|--------|----------|
| **P1** | ✅ shipped | create, add, ls; 17 codecs; Finder selection; interactive TTY picker; recents; dry-run; backup; slow-job notifications |
| **P2** | ✅ shipped | extract archives + disk images; test/verify integrity; clean-in-place delete; junk groups (macos/windows/vcs/hidden); destination picker; pre-flight verify |
| **P3** | planned | watch/folder mode with auto-compression; profiles/presets; man page; `brew tap` distribution; remote archive support

