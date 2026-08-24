# Task 18 Review: Completion notifications

## Spec compliance: ✅

- `internal/ui/notify.go` created with `Notify(title, body string)` and `NotifyIfSlow(start time.Time, d time.Duration, title, body string)` signatures matching brief exactly.
- `Notify` uses `osascript -e 'display notification ...'` with best-effort `_ = .Start()` — never fails the run.
- `NotifyIfSlow` gates on `time.Since(start) > d` and calls `Notify` when slow.
- File compiles cleanly (`go build ./internal/ui/`).

## Quality: Approved

Wiring into `runCreate` / `addCmd.RunE` is explicitly deferred per brief; no defect raised.

Commit: `1297b9d feat(ui): slow-job completion notifications`
