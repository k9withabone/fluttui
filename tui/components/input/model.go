package input

import (
	"errors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/constants"
)

type Options struct {
	Text            string
	Prompt          string
	Placeholder     string
	Footer          string
	Footer2         string
	SavePlaceholder bool
}

type Model struct {
	Text            string
	input           textinput.Model
	footer          string
	footer2         string
	help            help.Model
	savePlaceholder bool
}

func New(options Options) (Model, tea.Cmd, error) {
	input := textinput.New()
	input.SetValue(options.Text)
	input.Prompt = options.Prompt
	input.Placeholder = options.Placeholder

	input.Focus()

	m := Model{
		Text: options.Text,
		input: input,
		footer: options.Footer,
		footer2: options.Footer2,
		help: help.New(),
		savePlaceholder: options.SavePlaceholder,
	}

	// update to populate m.Text
	model, _ := m.Update(nil)
	newModel, ok := model.(Model)
	if !ok {
		return m, nil, errors.New("Error creating input")
	}
	m = newModel

	return m, textinput.Blink, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	m.Text = m.input.Value()

	if m.savePlaceholder && m.Text == "" {
		m.Text = m.input.Placeholder
	}

	return m, cmd
}

func (m Model) View() string {
	footer := m.footer
	if m.footer2 != "" {
		if m.Text != "" {
			footer += m.Text
		} else {
			footer += m.input.Placeholder
		}
		footer += m.footer2
	}

	return constants.ViewStyle.Render(
		m.input.View() + "\n\n" +
	    footer + "\n\n" +
	    m.help.View(constants.Keys),
	)
}