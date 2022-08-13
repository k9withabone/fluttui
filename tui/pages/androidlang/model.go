package androidlang

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/choosebutton"
)

const (
	kotlin = iota
	java
)

type Model struct {
	Choice string
	choosebutton choosebutton.Model
}

func New(androidLang string) Model {
	choice := kotlin
	if androidLang == "java" {
		choice = java
	}

	return Model{
		Choice: androidLang,
		choosebutton: choosebutton.New(choosebutton.Options{
			Choice: choice,
			Prompt: "Flutter Android Language?",
			Items: []string{
				"Kotlin",
				"Java",
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

	m.Choice = "kotlin"
	if m.choosebutton.Choice == java {
		m.Choice = "java"
	}

	return m, cmd
}

func (m Model) View() string {
	return m.choosebutton.View();
}