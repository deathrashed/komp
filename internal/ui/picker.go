package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	dirStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
)

type pickerModel struct {
	cwd      string
	entries  []os.DirEntry
	cursor   int
	offset   int
	chosen   map[string]bool
	order    []string
	quitting bool
	aborted  bool
	height   int
	width    int
	err      error
}

func newPickerModel(root string) pickerModel {
	m := pickerModel{
		cwd:    root,
		chosen: map[string]bool{},
		height: 20,
		width:  60,
	}
	_ = m.reload()
	return m
}

func (m *pickerModel) reload() error {
	entries, err := os.ReadDir(m.cwd)
	if err != nil {
		m.err = err
		m.entries = nil
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		ai, aj := entries[i].IsDir(), entries[j].IsDir()
		if ai != aj {
			return ai // dirs first
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})
	m.entries = entries
	if m.cursor >= len(m.entries) {
		m.cursor = max(0, len(m.entries)-1)
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.ensureVisible()
	return nil
}

func (m *pickerModel) ensureVisible() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *pickerModel) up() {
	parent := filepath.Dir(m.cwd)
	if parent != m.cwd {
		m.cwd = parent
		m.cursor = 0
		m.offset = 0
		_ = m.reload()
	}
}

func (m *pickerModel) into() {
	if len(m.entries) == 0 {
		return
	}
	if m.cursor < 0 || m.cursor >= len(m.entries) {
		return
	}
	e := m.entries[m.cursor]
	full := filepath.Join(m.cwd, e.Name())
	st, err := os.Stat(full)
	if err != nil {
		return
	}
	if st.IsDir() {
		m.cwd = full
		m.cursor = 0
		m.offset = 0
		_ = m.reload()
	}
}

func (m *pickerModel) toggle() {
	if len(m.entries) == 0 {
		return
	}
	if m.cursor < 0 || m.cursor >= len(m.entries) {
		return
	}
	e := m.entries[m.cursor]
	full := filepath.Join(m.cwd, e.Name())
	if m.chosen[full] {
		delete(m.chosen, full)
		// remove from order preserving remaining order
		filtered := m.order[:0]
		for _, p := range m.order {
			if p != full {
				filtered = append(filtered, p)
			}
		}
		m.order = filtered
	} else {
		m.chosen[full] = true
		m.order = append(m.order, full)
	}
}

func runPicker(root string) ([]string, error) {
	if root == "" {
		root, _ = os.Getwd()
	}
	abs, err := filepath.Abs(root)
	if err == nil {
		root = abs
	}
	if st, err := os.Stat(root); err != nil || !st.IsDir() {
		if cwd, err2 := os.Getwd(); err2 == nil {
			if abs2, err3 := filepath.Abs(cwd); err3 == nil {
				root = abs2
			} else {
				root = cwd
			}
		}
	}
	m := newPickerModel(root)
	p := tea.NewProgram(m)
	fm, err := p.Run()
	if err != nil {
		return nil, err
	}
	final := fm.(pickerModel)
	if final.aborted {
		return nil, errors.New("aborted")
	}
	return final.order, nil
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = max(10, msg.Height-4)
		m.ensureVisible()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			m.aborted = true
			return m, tea.Quit
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case " ":
			m.toggle()
			return m, nil
		case "h", "left", "backspace":
			m.up()
			return m, nil
		case "l", "right":
			m.into()
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.ensureVisible()
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				m.ensureVisible()
			}
			return m, nil
		case "g", "home":
			m.cursor = 0
			m.offset = 0
			return m, nil
		case "G", "end":
			if len(m.entries) > 0 {
				m.cursor = len(m.entries) - 1
				m.ensureVisible()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	if m.quitting {
		return ""
	}
	var b strings.Builder
	status := fmt.Sprintf("cwd: %s  selected: %d  [↑/k ↓/j] move [space]toggle [l/→]open [h/←]up [enter]done [esc]abort", m.cwd, len(m.chosen))
	b.WriteString(statusStyle.Render(status))
	b.WriteString("\n")
	if m.err != nil {
		b.WriteString(fmt.Sprintf("error reading %s: %v\n", m.cwd, m.err))
		return b.String()
	}
	if len(m.entries) == 0 {
		b.WriteString("(empty)\n")
		return b.String()
	}
	end := min(m.offset+m.height, len(m.entries))
	for i := m.offset; i < end; i++ {
		e := m.entries[i]
		full := filepath.Join(m.cwd, e.Name())
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		checked := "[ ] "
		if m.chosen[full] {
			checked = "[x] "
		}
		name := e.Name()
		if e.IsDir() {
			name = dirStyle.Render(name + "/")
		}
		line := cursor + checked + name
		if i == m.cursor {
			line = cursorStyle.Render(line)
		} else if m.chosen[full] {
			line = selectedStyle.Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	if len(m.entries) > m.height {
		b.WriteString(statusStyle.Render(fmt.Sprintf("— %d/%d —", m.cursor+1, len(m.entries))))
		b.WriteString("\n")
	}
	return b.String()
}

