package ioslang

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/choosebutton"
)

const (
	swift = iota
	objc
)

type Model struct {
	Choice string
	choosebutton choosebutton.Model
}

func New(iosLang string) Model {
	choice := swift
	if iosLang == "objc" {
		choice = objc
	}

	return Model{
		Choice: iosLang,
		choosebutton: choosebutton.New(choosebutton.Options{
			Choice: choice,
			Prompt: "Flutter iOS Language?",
			Items: []string{
				"Swift",
				"Objective-C",
			},
		}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.choosebutton, cmd = m.choosebutton.Update(msg)

	m.Choice = "swift"
	if m.choosebutton.Choice == objc {
		m.Choice = "objc"
	}

	return m, cmd
}

func (m Model) View() string {
	return m.choosebutton.View();
}