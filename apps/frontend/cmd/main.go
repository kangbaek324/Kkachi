package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kangbaek324/kkachi/apps/frontend/internal/app"
	"github.com/kangbaek324/kkachi/apps/frontend/internal/config"
)

func main() {
	cfg := config.Load()

	p := tea.NewProgram(app.NewModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running tui: %v\n", err)
		os.Exit(1)
	}
}
