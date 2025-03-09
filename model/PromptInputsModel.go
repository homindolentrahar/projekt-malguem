package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type PromptInputStringModel struct {
	Prompt string
	Input  string
	Done   bool
}

func (i PromptInputStringModel) Init() tea.Cmd {
	return nil
}

func (i PromptInputStringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i.Done = true
			return i, tea.Quit
		case "backspace":
			if len(i.Input) > 0 {
				i.Input = i.Input[:len(i.Input)-1]
			}
		case "ctrl+c", ":q":
			return i, tea.Quit
		default:
			i.Input += msg.String()
		}
	}

	return i, nil
}

func (i PromptInputStringModel) View() string {
	if i.Done {
		return "\n"
	}

	return fmt.Sprintf("%s: %s", i.Prompt, i.Input)
}
