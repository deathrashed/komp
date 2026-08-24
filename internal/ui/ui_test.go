package ui

import "testing"

func TestInteractiveFalseWhenStdoutPiped(t *testing.T) {
	// In `go test` stdout is captured (non-tty).
	if Interactive() {
		t.Skip("stdout is a terminal in this environment")
	}
}
