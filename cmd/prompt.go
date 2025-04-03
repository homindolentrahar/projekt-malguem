package cmd

import (
	"malguem/model"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func PrompInput(prompt string) string {
	p := tea.NewProgram(model.PrompInput{Prompt: prompt})
	stringInput, err := p.Run()

	if err != nil {
		os.Exit(1)
	}

	value := stringInput.(model.PrompInput).Input

	if value == "" {
		return ""
	}

	return value
}

func PromptChoice(prompt string, options []string) string {
	p := tea.NewProgram(model.PromptChoice{Prompt: prompt, Options: options})
	selected, err := p.Run()

	if err != nil {
		os.Exit(1)
	}

	value := selected.(model.PromptChoice).Selected

	if value == "" {
		return ""
	}

	return value
}
