package main

import (
	"log"

	"github.com/IvanYaremko/rssdukester/models/container"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(container.Container{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
