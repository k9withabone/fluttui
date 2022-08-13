package directory

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/input"
	"github.com/k9withabone/fluttui/tui/constants"
)

func New(directory string, name string) (input.Model, tea.Cmd, error) {
	placeholder := "./" + name

	return input.New(input.Options{
		Text: directory,
		Prompt: constants.PromptStyle.Render("Flutter Project Directory? (Optional)") + "\n\n> ",
		Placeholder: placeholder,
		Footer: "If left blank \"" + placeholder + "\" will be used.",
	})
}