package main

import (
	"log"

	"github.com/asahnoln/go-planner/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := tea.NewProgram(tui.New()).Start(); err != nil {
		log.Fatalf("could not start the program: %q\n", err)
	}
}
