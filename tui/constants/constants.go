package constants

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type CompletedMsg struct{
	Err    error
	Output string
}

type keyMap struct {
	Back key.Binding
	Next key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Next,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.Back,
		k.Next,
		k.Quit,
	}}
}

var (
	Keys = keyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Next: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "next"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+d"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}

	White       = lipgloss.Color("230")

	FlutterNavy = lipgloss.Color("#042B59")
	FlutterBlue = lipgloss.Color("#0553B1")
	FlutterSky  = lipgloss.Color("#027DFD")

	PromptStyle = lipgloss.NewStyle().
                  Background(FlutterBlue).
                  Foreground(White).
                  Padding(0, 1)

	ViewStyle   = lipgloss.NewStyle().Padding(1, 0, 0, 2)
)