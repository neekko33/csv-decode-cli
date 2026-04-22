package app

import "github.com/charmbracelet/bubbles/textinput"

type step int

const (
	stepInputCSV step = iota
	stepSelectFields
	stepOutputPath
	stepConflict
	stepDone
)

type model struct {
	step step

	input textinput.Model

	csvPath     string
	outputPath  string
	headers     []string
	selected    map[int]bool
	cursor      int
	conflictPos int

	errorMsg string
	doneMsg  string
}

func NewModel() model {
	ti := textinput.New()
	ti.Placeholder = "Path to input CSV"
	ti.Prompt = "> "
	ti.Focus()

	return model{
		step:     stepInputCSV,
		input:    ti,
		selected: map[int]bool{},
	}
}
