package overwrite

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/choosebutton"
)

const (
	no = iota
	yes
)

type Model struct {
	Overwrite    bool
	choosebutton choosebutton.Model
}

func New(overwrite bool) Model {
	choice := no
	if overwrite {
		choice = yes
	}

	return Model{
		Overwrite: overwrite,
		choosebutton: choosebutton.New(choosebutton.Options{
			Choice: choice,
			Prompt: "Overwrite?",
			Items: []string{
				"No",
				"Yes",
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

	m.Overwrite = m.choosebutton.Choice == yes

	return m, cmd
}

func (m Model) View() string {
	return m.choosebutton.View();
}