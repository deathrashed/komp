package ui

import (
	"os/exec"
	"time"
)

// Notify posts a macOS notification (best-effort, never fails the run).
func Notify(title, body string) {
	_ = exec.Command("osascript", "-e",
		`display notification "`+body+`" with title "`+title+`"`).Start()
}

// NotifyIfSlow fires Notify only when work took longer than d.
func NotifyIfSlow(start time.Time, d time.Duration, title, body string) {
	if time.Since(start) > d {
		Notify(title, body)
	}
}
