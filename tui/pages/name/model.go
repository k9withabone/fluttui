package name

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/input"
	"github.com/k9withabone/fluttui/tui/constants"
)

func New(name string) (input.Model, tea.Cmd, error) {
	const placeholder = "flutter_project_name"

	if name == placeholder {
		name = ""
	}

	return input.New(input.Options{
		Text: name,
		Prompt: constants.PromptStyle.Render("Flutter Project Name?") + "\n\n> ",
		Placeholder: placeholder,
		Footer: "If left blank \"" + placeholder + "\" will be used.",
		SavePlaceholder: true,
	})
}