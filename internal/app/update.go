package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"csv-decode-cli/internal/csvsvc"
)

func (m model) Init() tea.Cmd {
	return textBlinkCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	switch m.step {
	case stepInputCSV:
		return m.updateInputCSV(msg)
	case stepSelectFields:
		return m.updateSelectFields(msg)
	case stepOutputPath:
		return m.updateOutputPath(msg)
	case stepConflict:
		return m.updateConflict(msg)
	case stepDone:
		return m.updateDone(msg)
	default:
		return m, nil
	}
}

func setHomeDir(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home
	}

	if strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~"+string(filepath.Separator)))
	}

	return path
}

func (m model) updateInputCSV(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "tab":
			if completed, changed := completePath(m.input.Value()); changed {
				m.input.SetValue(completed)
				m.input.CursorEnd()
			}
			return m, nil
		case "enter":
			path := setHomeDir(strings.TrimSpace(m.input.Value()))
			headers, err := csvsvc.ReadHeaders(path)
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}

			m.csvPath = path
			m.headers = headers
			m.selected = map[int]bool{}
			m.cursor = 0
			m.step = stepSelectFields
			m.errorMsg = ""
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateSelectFields(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < len(m.headers)-1 {
			m.cursor++
		}
	case " ":
		if len(m.headers) > 0 {
			m.selected[m.cursor] = !m.selected[m.cursor]
		}
	case "enter":
		fields := m.selectedFields()
		if len(fields) == 0 {
			m.errorMsg = "choose at least one field"
			return m, nil
		}

		m.step = stepOutputPath
		m.errorMsg = ""
		m.input = newOutputInput(csvsvc.DefaultOutputPath(m.csvPath))
		return m, nil
	}

	return m, nil
}

func (m model) updateOutputPath(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "tab":
			if completed, changed := completePath(m.input.Value()); changed {
				m.input.SetValue(completed)
				m.input.CursorEnd()
			}
			return m, nil
		case "enter":
			outputPath := setHomeDir(strings.TrimSpace(m.input.Value()))
			m.outputPath = outputPath

			err := csvsvc.ValidateDestination(outputPath, false)
			if err != nil {
				if errors.Is(err, csvsvc.ErrDestinationExists) {
					m.step = stepConflict
					m.conflictPos = 0
					m.errorMsg = ""
					return m, nil
				}
				m.errorMsg = err.Error()
				return m, nil
			}

			return m.runDecode(false)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateConflict(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up", "down":
		if m.conflictPos == 0 {
			m.conflictPos = 1
		} else {
			m.conflictPos = 0
		}
	case "enter":
		if m.conflictPos == 0 {
			return m.runDecode(true)
		}
		m.step = stepOutputPath
		m.input = newOutputInput(m.outputPath)
		return m, nil
	}

	return m, nil
}

func (m model) updateDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) runDecode(allowOverwrite bool) (tea.Model, tea.Cmd) {
	err := csvsvc.DecodeCSVFields(m.csvPath, m.outputPath, m.selectedFields(), allowOverwrite)
	if err != nil {
		if errors.Is(err, csvsvc.ErrDestinationExists) {
			m.step = stepConflict
			m.errorMsg = ""
			return m, nil
		}
		m.errorMsg = err.Error()
		m.step = stepOutputPath
		m.input = newOutputInput(m.outputPath)
		return m, nil
	}

	m.doneMsg = fmt.Sprintf("Decode complete.\nOutput: %s\nFields: %s", m.outputPath, strings.Join(m.selectedFields(), ", "))
	m.errorMsg = ""
	m.step = stepDone
	return m, nil
}

func (m model) selectedFields() []string {
	fields := make([]string, 0, len(m.selected))
	for idx, checked := range m.selected {
		if checked && idx >= 0 && idx < len(m.headers) {
			fields = append(fields, m.headers[idx])
		}
	}
	// Keep header order stable.
	ordered := make([]string, 0, len(fields))
	for i := range m.headers {
		if m.selected[i] {
			ordered = append(ordered, m.headers[i])
		}
	}
	return ordered
}

func textBlinkCmd() tea.Cmd {
	return textinput.Blink
}

func newOutputInput(defaultPath string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.SetValue(defaultPath)
	ti.CursorEnd()
	ti.Focus()
	return ti
}
