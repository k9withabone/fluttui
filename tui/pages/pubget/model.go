package pubget

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/choosebutton"
)

const (
	PubGetYes = iota
	PubGetOffline
	PubGetNo
)

type Model struct {
	Choice       int
	choosebutton choosebutton.Model
}

func New(choice int) Model {
	return Model{
		Choice: choice,
		choosebutton: choosebutton.New(choosebutton.Options{
			Choice: choice,
			Prompt: "flutter pub get?",
			Items: []string{
				"Yes",
				"Offline",
				"No",
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

	m.Choice = m.choosebutton.Choice

	return m, cmd
}

func (m Model) View() string {
	return m.choosebutton.View();
}