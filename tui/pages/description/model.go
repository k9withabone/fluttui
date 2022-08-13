package description

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/input"
	"github.com/k9withabone/fluttui/tui/constants"
)

func New(description string) (input.Model, tea.Cmd, error) {
	return input.New(input.Options{
		Text: description,
		Prompt: constants.PromptStyle.Render("Flutter Project Description? (Optional)") + "\n\n> ",
		Placeholder: "Description...",
		Footer: "Defaults to \"A new Flutter project.\"",
	})
}