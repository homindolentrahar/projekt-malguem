package cmd

import (
	"malguem/model"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func PrompInput(prompt string) string {
	p := tea.NewProgram(model.PromptInputStringModel{Prompt: prompt})
	stringInput, err := p.Run()

	if err != nil {
		os.Exit(1)
	}

	value := stringInput.(model.PromptInputStringModel).Input

	if value == "" {
		return "default"
	}

	return value
}
