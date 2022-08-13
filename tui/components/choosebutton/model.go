package choosebutton

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/k9withabone/fluttui/tui/constants"
)

var (
	unselectedStyle = lipgloss.NewStyle().
	                  Background(lipgloss.AdaptiveColor{Dark: "235", Light: "252"}).
	                  Padding(0, 1).
					  Margin(1, 2, 1, 0)
	selectedStyle   = unselectedStyle.Copy().
	                  Background(constants.FlutterBlue).
					  Foreground(constants.White).
					  Bold(true)
)

type keyMap struct{
	left  key.Binding
	right key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return append(
		[]key.Binding{
			k.left,
			k.right,
		},
		constants.Keys.ShortHelp()...,
	)
}

func (k keyMap) FullHelp() [][]key.Binding {
	return append(
		[][]key.Binding{{
			k.left,
			k.right,
		}},
		constants.Keys.FullHelp()...,
	)
}

type Options struct {
	Choice      int
	Prompt      string
	Items       []string
	DisableHelp bool
	Center      bool
}

type Model struct {
	Choice      int
	prompt      string
	items       []string
	help        help.Model
	KeyMap      keyMap
	disableHelp bool
	center      bool
}

func New(options Options) Model {
	return Model{
		Choice: options.Choice,
		prompt: options.Prompt,
		items: options.Items,
		help: help.New(),
		KeyMap: keyMap{
			left: key.NewBinding(
				key.WithKeys("left", "h"),
				key.WithHelp("←/h", "left"),
			),
			right: key.NewBinding(
				key.WithKeys("right", "l"),
				key.WithHelp("→/l", "right"),
			),
		},
		disableHelp: options.DisableHelp,
		center: options.Center,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.left):
			m.Choice--
			if m.Choice < 0 {
				m.Choice = len(m.items) - 1
			}
		
		case key.Matches(msg, m.KeyMap.right):
			m.Choice++
			if m.Choice >= len(m.items) {
				m.Choice = 0
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var itemStrings []string
	for i, item := range m.items {
		if i == m.Choice {
			itemStrings = append(itemStrings, selectedStyle.Render(item))
		} else {
			itemStrings = append(itemStrings, unselectedStyle.Render(item))
		}
	}
	
	strings := []string{
		constants.PromptStyle.Render(m.prompt),
		lipgloss.JoinHorizontal(lipgloss.Top, itemStrings...),
	}
	if !m.disableHelp {
		strings = append(strings, m.help.View(m.KeyMap))
	}

	position := lipgloss.Left
	if m.center {
		position = lipgloss.Center
	}

	return constants.ViewStyle.Render(
		lipgloss.JoinVertical(position, strings...),
	)
}