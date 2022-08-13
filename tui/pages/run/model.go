package run

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"github.com/k9withabone/fluttui/tui/constants"
)

var commandStyle = lipgloss.NewStyle().Margin(1, 0)

type Model struct {
	title     string
	spinner   spinner.Model
	command   string
	completed bool
	output    string
}

func New(command string) (Model, tea.Cmd) {
	// put quotes around description
	const desc = "--description"
	if strings.Contains(command, desc) {
		cmdStrings := strings.Split(command, " ")

		var descStart, descEnd int
		for i, cmdStr := range cmdStrings {
			if cmdStr == desc {
				descStart = i + 1
			} else if descStart != 0 && cmdStr[0] == '-' {
				descEnd = i
				break
			}
		}
		
		if descEnd == 0 {
			descEnd = len(cmdStrings)
		}

		command =  strings.Join(cmdStrings[:descStart], " ")
		command += " \"" + strings.Join(cmdStrings[descStart:descEnd], " ") + "\" "
		command += strings.Join(cmdStrings[descEnd:], " ")
	}

	m := Model{
		title: "Running",
		spinner: spinner.New(
			spinner.WithStyle(lipgloss.NewStyle().Foreground(constants.FlutterSky)),
		),
		command: command,
	}

	return m, m.spinner.Tick
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case constants.CompletedMsg:
		m.completed = true
		m.output = msg.Output
		m.title = "Finished!"
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	title := constants.PromptStyle.Render(m.title)
	if !m.completed {
		title += " " + m.spinner.View()
	}

	command := wordwrap.String(m.command, 80)

	str := lipgloss.JoinVertical(lipgloss.Left,
		title,
		commandStyle.Render(command),
	)

	if m.completed {
		str += "\nOutput:\n\n" + m.output
	}

	return constants.ViewStyle.Render(str)
}