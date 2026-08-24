<div align="center">
  <img src="assets/icon.png" alt="komp" width="200">

  <h1>KOMP</h1>

  <p><strong>One binary for every archive on macOS — compress, extract, peek, clean, test, convert, and build disk images.</strong></p>

  <p>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/go-1.27+-1e1e1e?style=for-the-badge&logo=go&logoColor=01acd7" alt="Go 1.27+"></a>
    <a href="https://github.com/deathrashed/komp/releases"><img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux-1e1e1e?style=for-the-badge&logo=apple&logoColor=01acd7" alt="Platform"></a>
    <a href="https://github.com/deathrashed/komp/releases"><img src="https://img.shields.io/badge/version-v1.0.0-1e1e1e?style=for-the-badge&logo=github&logoColor=01acd7" alt="Version"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-1e1e1e?style=for-the-badge&logo=openaccess&logoColor=01acd7" alt="MIT License"></a>
  </p>

  <p>
    <a href="#quick-start">Quick Start</a> |
    <a href="#usage">Usage</a> |
    <a href="#commands">Commands</a> |
    <a href="#interactive-mode">Interactive Mode</a> |
    <a href="#codecs">Codecs</a> |
    <a href="#junk-cleaning">Junk Cleaning</a> |
    <a href="#disk-images">Disk Images</a> |
    <a href="#building">Building</a>
  </p>
</div>

---

## <img src="https://api.iconify.design/mdi:rocket-launch-outline.svg?color=%2301acd7" height="22"> Quick Start

**Install a prebuilt binary** (macOS Apple Silicon):

```bash
curl -fsSL -o komp https://github.com/deathrashed/komp/releases/latest/download/komp-darwin-arm64.tar.gz
tar xzf komp-darwin-arm64.tar.gz
sudo mv komp /usr/local/bin/
```

**Or build from source:**

```bash
go install github.com/deathrashed/komp@latest
```

> [!NOTE]
> `komp` collides with nothing on macOS — no system tool or common framework uses the name, so it is safe in `$PATH`.

## <img src="https://api.iconify.design/mdi:console.svg?color=%2301acd7" height="22"> Usage

```bash
komp                          # interactive menu — pick an action
komp -f zstd notes.txt        # compress to zstd
komp -f zip --delete ~/proj   # compress folder, remove originals
komp add project.zip new.go   # add files to an existing archive
komp ls backup.tar.gz         # peek inside
komp un photos.zip            # extract
komp clean --groups macos,vcs app.zip     # strip junk members
komp t backup.7z              # verify integrity
komp cv old.tar.gz --to zst   # convert tar.gz → tar.zst (streamed)
komp img ~/disk --type dmg    # build a disk image
```

Run `komp` with **no arguments** in a terminal for a guided, multi-page menu — it filters archive lists to only the formats each action supports, keeps your previous answers on screen, and `esc` walks you back a page.

## <img src="https://api.iconify.design/mdi:gesture-tap-button.svg?color=%2301acd7" height="22"> Commands

| Command | What it does |
| --- | --- |
| `komp [format] [files...]` | Compress files or folders |
| `komp add <archive> [files...]` | Add files into an existing archive |
| `komp ls <archive>` | List archive contents |
| `komp un <archive>...` | Extract archives and disk images |
| `komp clean <archive>...` | Strip junk members in place |
| `komp t <archive>...` | Test archive integrity |
| `komp cv <archive>` | Recompress to another format |
| `komp img <folder>` | Build dmg / sparseimage / iso / pkg |
| `komp completions <shell>` | Generate shell completions |
| `komp man` | Print a man page |

### <img src="https://api.iconify.design/mdi:tune.svg?color=%2301acd7" height="18"> Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| `--format` | `-f` | — | codec name (also accepted as first bare arg) |
| `--delete` | `-d` | off | remove originals after success |
| `--output` | `-o` | input dir | output directory |
| `--level` | `-L` | `normal` | `fast` \| `normal` \| `max` |
| `--separate` | | off | one archive per input |
| `--each` | | off | streams: compress each input separately |
| `--finder` | | off | use current Finder selection as input |
| `--dry-run` | | off | print plan, touch nothing |
| `--backup` | | off | back up overwritten targets |

## <img src="https://api.iconify.design/mdi:form-select.svg?color=%2301acd7" height="22"> Interactive Mode

`komp` with no arguments opens a `huh`-powered flow:

- **Select Action** — Compress, Add, Peek, Extract, Clean, Test, Convert, Build, Quit
- **Filtered archive lists** — Add only lists formats that support append; Clean only lists cleanable ones; Convert only extractable sources; Extract includes `.dmg` / `.iso` / `.sparseimage` / `.sparsebundle`
- **Inline multi-page forms** — previous selections stay visible above the current field; `shift+tab` steps back through fields
- **esc** — returns to the previous menu, or quits from the top
- **Spinner** — long jobs show an inline progress spinner and a macOS notification on completion

## <img src="https://api.iconify.design/mdi:folder-zip-outline.svg?color=%2301acd7" height="22"> Codecs

17 formats out of the box. Missing binaries are detected and flagged with the Homebrew formula to install.

| Built-in (system) | Via Homebrew |
| --- | --- |
| zip, tar, gzip, bzip2, aar (`aa`), dmg/iso (`hdiutil`) | `sevenzip` (7z), `rar`, `zpaq`, `xz`, `lzip`, `zstd`, `lz4`, `brotli`, `lzop`, `snzip`, `lrzip` |

<details>
<summary><strong>Codec capability matrix</strong></summary>

| Codec | Create | Add | List | Test | Extract | Clean in place |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| zip | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 7z | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| tar (+ gz/bz2/xz/zst/br/lz4/lrz wrappers) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ (rebuild) |
| rar | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| zpaq | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| aar (`aa`) | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| streams (gz, bz2, xz, lzma, lz, zst, lz4, br, lzo, sz, lrz) | ✓ | — | — | ✓ | ✓ | — |

</details>

## <img src="https://api.iconify.design/mdi:broom.svg?color=%2301acd7" height="22"> Junk Cleaning

`komp clean` strips known junk members from archives without touching the rest:

| Group | Matches |
| --- | --- |
| `macos` | `.DS_Store`, `._*`, `__MACOSX/`, `.Spotlight-V100`, `.Trashes`, `.fseventsd`, `.LSOverride`, `Icon\r` |
| `windows` | `Thumbs.db`, `desktop.ini`, `$RECYCLE.BIN/` |
| `vcs` | `.git/`, `.svn/`, `.hg/` (any depth) |
| `hidden` | dotfiles not claimed by another group |

zip and 7z are cleaned in place under an atomic swap; tar-family archives are rebuilt and swapped. `--dry-run` previews, `--backup` keeps `<archive>.bak`.

## <img src="https://api.iconify.design/mdi:disc.svg?color=%2301acd7" height="22"> Disk Images

```bash
komp img ~/app --type dmg --volname "My App"       # compressed UDZO dmg
komp img ~/vault --type sparseimage --size 4g      # APFS sparse image
komp img ~/music --type iso                        # hybrid ISO
komp img ~/pkgroot --type pkg --id com.me.tool     # macOS installer pkg
```

| Type | Tool | Notes |
| --- | --- | --- |
| `dmg` | `hdiutil` | UDZO compressed by default |
| `sparsebundle` / `sparseimage` | `hdiutil` | requires `--size`, APFS default |
| `iso` | `hdiutil makehybrid` | ISO + Joliet |
| `pkg` | `pkgbuild` | requires `--id`, auto-version `1.0` |

## <img src="https://api.iconify.design/mdi:exit-to-app.svg?color=%2301acd7" height="22"> Exit Codes

| Code | Meaning |
| --- | --- |
| 0 | success |
| 1 | runtime / environment error (no input, aborted, osascript failure) |
| 2 | usage error (unknown format, missing binary) |
| 3 | execution error (compressor / extractor failed) |
| 4 | integrity failure |

## <img src="https://api.iconify.design/mdi:file-tree-outline.svg?color=%2301acd7" height="22"> Project Structure

```text
komp/
├── main.go                    # entry point
├── assets/icon.png            # header image
├── internal/
│   ├── cli/                   # cobra commands + interactive flow
│   ├── codec/                 # 17-codec capability table
│   ├── engine/                # create/add/ls/un/clean/test/convert pipelines
│   ├── image/                 # dmg/iso/pkg/sfx builders
│   ├── junk/                  # junk member patterns + matcher
│   ├── recents/               # recent-archives store
│   ├── selection/             # Finder selection via osascript
│   ├── ui/                    # huh forms, file picker, spinner, notifications
│   └── xdg/                   # XDG config/data paths
└── CHANGELOG.md               # release history
```

## <img src="https://api.iconify.design/mdi:hammer-wrench.svg?color=%2301acd7" height="22"> Building

```bash
git clone https://github.com/deathrashed/komp
cd komp
go build -o komp .
go test ./...
```

<details>
<summary><strong>Cross-compile all release targets</strong></summary>

```bash
mkdir -p dist
for t in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64; do
  os=${t%/*}; arch=${t#*/}
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -trimpath -ldflags "-s -w" -o dist/komp .
  tar czf dist/komp-$os-$arch.tar.gz -C dist komp
done
```

</details>

## <img src="https://api.iconify.design/mdi:license.svg?color=%2301acd7" height="22"> License

[MIT](LICENSE) — see [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.
