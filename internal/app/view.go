package app

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	switch m.step {
	case stepInputCSV:
		return m.viewInputCSV()
	case stepSelectFields:
		return m.viewSelectFields()
	case stepOutputPath:
		return m.viewOutputPath()
	case stepConflict:
		return m.viewConflict()
	case stepDone:
		return m.viewDone()
	default:
		return ""
	}
}

func (m model) viewInputCSV() string {
	s := "CSV Decode CLI\n\nEnter input CSV file path:\n"
	s += m.input.View()
	s += "\n\nPress Tab to autocomplete path."
	s += "\nPress Enter to continue, q to quit."
	if m.errorMsg != "" {
		s += "\n\nError: " + m.errorMsg
	}
	return s + "\n"
}

func (m model) viewSelectFields() string {
	var b strings.Builder
	b.WriteString("Select fields to decode\n\n")
	b.WriteString("Use up/down, space to toggle, enter to confirm.\n\n")

	for i, h := range m.headers {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		checked := " "
		if m.selected[i] {
			checked = "x"
		}
		b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, h))
	}
	if m.errorMsg != "" {
		b.WriteString("\nError: " + m.errorMsg + "\n")
	}
	b.WriteString("\nPress q to quit.\n")
	return b.String()
}

func (m model) viewOutputPath() string {
	s := "Enter output CSV file path:\n"
	s += m.input.View()
	s += "\n\nDefault path uses input directory with '-decoded' suffix."
	s += "\nPress Tab to autocomplete path."
	s += "\nPress Enter to continue, q to quit."
	if m.errorMsg != "" {
		s += "\n\nError: " + m.errorMsg
	}
	return s + "\n"
}

func (m model) viewConflict() string {
	options := []string{"Overwrite existing file", "Enter a different output path"}

	var b strings.Builder
	b.WriteString("Output file already exists:\n")
	b.WriteString(m.outputPath + "\n\n")
	b.WriteString("Choose an action:\n\n")
	for i, op := range options {
		cursor := " "
		if i == m.conflictPos {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, op))
	}
	b.WriteString("\nUse up/down and Enter.\n")
	return b.String()
}

func (m model) viewDone() string {
	return m.doneMsg + "\n\nPress Enter to exit.\n"
}
