# Contributing to komp

Thanks for your interest in contributing!

## Getting started

```bash
git clone https://github.com/deathrashed/komp
cd komp
go build -o komp .
go test ./...
```

## Ground rules

- **Tests first for engine changes.** Anything touching `internal/engine/` or `internal/codec/` needs table-driven tests. Run `go test ./...` before opening a PR.
- **Keep the binary dependency-light.** Runtime dependencies are external CLI tools (`zip`, `7z`, `tar`, `hdiutil`...), never Go packages. New codecs are table entries in `internal/codec/table.go` with capability flags, not new code paths.
- **Respect the exit codes.** See the table in the README. New error classes go through `cli.classify`.
- **Interactive flows must degrade gracefully.** Every TUI path (`internal/ui/`) needs a non-TTY fallback that works when stdout is piped.
- **No comments unless the code is genuinely non-obvious.** The codec table's placeholder syntax (`{in}`, `{out}`, `{level}`, `{indir}`, `{inbase}`, `{dest}`) is documented in `internal/codec/types.go`.

## Adding a codec

1. Add an entry to `internal/codec/table.go` with `CreateArgs`, `TestArgs`, `ListArgs`, and — if supported — `AddArgs`, `ExtractArgs`, `DeleteArgs`.
2. Empty capability slices mean "unsupported" — the interactive pickers filter on these automatically.
3. Add the Homebrew formula name so `komp` can hint at installation.
4. Add a test in `internal/codec/` and, if the tool is installed, an engine round-trip test.

## Adding an image builder

Image builders live in `internal/image/table.go` using the same `{placeholder}` substitution. Validate required vars (size, id, ...) in `image.Build`.

## Commit style

Conventional commits: `feat(ui): ...`, `fix(engine): ...`, `docs: ...`.

## Reporting bugs

Open an issue with the exact command, expected vs actual output, and `komp`'s stderr. For extraction failures, include the archive format (not the archive itself).
