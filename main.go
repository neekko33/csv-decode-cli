package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"csv-decode-cli/internal/app"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		printUsage()
		return
	}

	if err := runApp(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  ./csv-decode-cli")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help    Show this help message and exit")
	fmt.Println("")
	fmt.Println("This command launches an interactive TUI.")
}

func runApp() error {
	p := tea.NewProgram(app.NewModel())
	_, err := p.Run()
	return err
}
