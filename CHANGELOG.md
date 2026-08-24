# Changelog

All notable changes to komp are documented here. Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-08-25

First public release.

### Added
- **Compress** — 17 codecs (zip, 7z, tar, rar, zpaq, aar, gzip, bzip2, xz, lzip, lzma, zstd, lz4, brotli, lzo, snappy, lrzip) with fast/normal/max levels, separate/each modes, dry-run, backup, and atomic output swaps.
- **Add** — append files to existing archives, including tar-family rebuilds, with `--last` / `--recent` targeting.
- **Peek** — list archive contents, including recent-archive shortcuts.
- **Extract** — no-clobber extraction for all archive codecs plus dmg / sparseimage / sparsebundle via attach-copy-detach, with destination picker.
- **Clean** — strip junk members (macOS, Windows, VCS, hidden groups) in place: zip/7z atomic swap, tar-family rebuild-and-swap, dry-run and backup support.
- **Test** — integrity verification with pre-flight hooks.
- **Convert** — recompress any extractable archive to another format, with streamed inner-tar fast path (tar.gz → tar.zst without full extraction).
- **Build** — dmg, sparsebundle, sparseimage, iso, and pkg builders via hdiutil/pkgbuild.
- **Interactive mode** — multi-page huh forms with capability-filtered archive lists, esc-back navigation, inline spinner, and macOS completion notifications.
- **Finder integration** — `--finder` operates on the current Finder selection.
- **Recents store** — XDG-backed recent-archive memory shared across commands.
- Shell completions (bash/zsh/fish/powershell) and a man page.

[1.0.0]: https://github.com/deathrashed/komp/releases/tag/v1.0.0
