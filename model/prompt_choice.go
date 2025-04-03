package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type PromptChoice struct {
	Cursor   int
	Prompt   string
	Options  []string
	Selected string
}

func (i PromptChoice) Init() tea.Cmd {
	return nil
}

func (i PromptChoice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return i, tea.Quit
		case "up", "k":
			if i.Cursor > 0 {
				i.Cursor--
			}
		case "down", "j":
			if i.Cursor < len(i.Options)-1 {
				i.Cursor++
			}
		case "enter":
			i.Selected = i.Options[i.Cursor]
			return i, tea.Quit
		}
	}

	return i, nil
}

func (i PromptChoice) View() string {
	prompt := i.Prompt + "\n"

	for index, option := range i.Options {
		cursor := ""

		if index == i.Cursor {
			cursor = ">"
		}

		prompt += fmt.Sprintf("%s %s\n", cursor, option)
	}

	prompt += "\nUse ↑/↓ to navigate, Enter to select, and Q to quit."

	if i.Selected != "" {
		prompt += "\nSelected: " + i.Selected
	}

	return prompt
}
