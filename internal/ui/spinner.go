package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type doneMsg struct{ err error }

type spinnerModel struct {
	title    string
	frames   []string
	i        int
	done     chan error
	err      error
	quitting bool
}

// RunWithSpinner runs fn while showing an inline spinner, then clears itself.
func RunWithSpinner(title string, fn func() error) error {
	if !Interactive() {
		return fn()
	}
	m := spinnerModel{
		title:  title,
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		done:   make(chan error, 1),
	}
	go func() { m.done <- fn() }()
	p := tea.NewProgram(m)
	fm, err := p.Run()
	if err != nil {
		return err
	}
	if sm, ok := fm.(spinnerModel); ok && sm.err != nil {
		return sm.err
	}
	return nil
}

func (m spinnerModel) Init() tea.Cmd { return tickCmd() }

func tickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg { return tickMsg{} })
}

type tickMsg struct{}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case doneMsg:
		m.err = msg.err
		m.quitting = true
		return m, tea.Quit
	case tickMsg:
		m.i = (m.i + 1) % len(m.frames)
		return m, tickCmd()
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("%s %s", m.frames[m.i], m.title)
}
