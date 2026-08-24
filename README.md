# komp

Unified archive & image toolkit for macOS. Compress, add, list, extract, clean, and convert archives; build dmg/iso/pkg images. One binary, 17 codecs, interactive TTY picker, and native Finder selection.

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

# Interactive mode — omit files/format in a terminal to get a TUI picker
komp
```

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

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | success |
| 1 | runtime / environment error (no input, aborted, osascript failure) |
| 2 | usage error (unknown format, missing binary, no format in non-interactive mode) |
| 3 | execution error (compressor / adder / lister failed) |

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
| Batch compress (one per file) | `komp --finder -f zip --separate --output ~/Archives` | Creates `~/Archives/<name>.zip` per selection |
| Dry-run plan | `komp --finder -f 7z --dry-run` | Prints plan to stderr, touches nothing |

## Apple-shadowing

`komp` collides with nothing on macOS. No system tool, built-in command, or common framework uses the name `komp`. Safe to drop into `$PATH` without shadowing.

## Roadmap

- **P1** (current) — create, add, ls; 17 codecs; Finder selection; interactive TTY picker; recents; dry-run; backup; slow-job notifications.
- **P2** — extract; test/verify integrity; clean-in-place delete; image building (dmg/iso/pkg); parallel `--separate`; shell completions.
- **P3** — watch/folder mode with auto-compression; profiles/presets; man page; `brew tap` distribution; remote archive support.
