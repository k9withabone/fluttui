package org

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/input"
	"github.com/k9withabone/fluttui/tui/constants"
)

func New(org string, name string) (input.Model, tea.Cmd, error) {
	const placeholder = "com.example"
	postfix := "." + name

	return input.New(input.Options{
		Text: org,
		Prompt: constants.PromptStyle.Render("Flutter Project Orginzation? (Optional)") + "\n\n> ",
		Placeholder: placeholder,
		Footer: "\"",
		Footer2: postfix + "\"\nwill be used in Java package names and in the iOS bundle identifier.",
	})
}