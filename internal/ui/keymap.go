package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// keymap binds esc (and ctrl+c) to abort so esc goes back/quits everywhere.
var keymap = func() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("esc", "ctrl+c"))
	return km
}()

// newForm builds a huh form with help and the shared keymap applied.
func newForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithShowHelp(true).WithKeyMap(keymap)
}
