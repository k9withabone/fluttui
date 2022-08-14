package confirm

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/k9withabone/fluttui/tui/components/choose"
	"github.com/k9withabone/fluttui/tui/components/choosebutton"
	"github.com/k9withabone/fluttui/tui/constants"
	"github.com/k9withabone/fluttui/tui/pages/pubget"
)

var (
	selectedStyle   = lipgloss.NewStyle().
	                  PaddingRight(2).
					  Border(lipgloss.RoundedBorder(), true).
					  BorderForeground(lipgloss.AdaptiveColor{
						Light: string(constants.FlutterBlue),
						Dark: string(constants.FlutterSky),
					  }).
	                  Margin(1, 1, 1, 0)
	unselectedStyle = selectedStyle.Copy().
	                  BorderForeground(lipgloss.AdaptiveColor{
						Light: "#C4DFFF",
						Dark: string(constants.FlutterNavy),
					  }).
	                  Foreground(lipgloss.AdaptiveColor{Light: "242", Dark: "245"})
)

type keyMap struct {
	optionsSelected         bool
	additionalShortHelpKeys func(bool) []key.Binding
	switchSelected          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return append(
		[]key.Binding{k.switchSelected},
		k.additionalShortHelpKeys(k.optionsSelected)...,
	)
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.switchSelected},
		k.additionalShortHelpKeys(k.optionsSelected),
	}
}

type Options struct {
	Template    string
	Name        string
	Directory   string
	Description string
	Org         string
	Platforms   []string
	Overwrite   bool
	PubGet      int
	IosLang     string
	AndroidLang string

	OptionsSelected bool
	ChooseSelected  int
	ConfirmChoice   int
}

type Model struct {
	choose       choose.Model
	choosebutton choosebutton.Model
	help         help.Model
	keyMap       keyMap

	OptionsSelected bool
	ChooseSelected  int
	ConfirmChoice   int
}

func New(options Options) (Model, error) {
	itemLabels := []string{
		"Template:",
		"Name:",
		"Directory:",
		"Description:",
		"Orginization:",
		"Platforms:",
		"Overwrite:",
		"pub get:",
		"iOS Language:",
		"Android Language: ",
	}

	itemLabelsBlock := lipgloss.JoinVertical(lipgloss.Left, itemLabels...)

	overwriteOption := "No"
	if options.Overwrite {
		overwriteOption = "Yes"
	}

	pubGetOption := "Yes"
	switch options.PubGet {
	case pubget.PubGetOffline:
		pubGetOption = "Offline"
	case pubget.PubGetNo:
		pubGetOption = "No"
	}

	itemOptions := []string{
		options.Template,
		options.Name,
		options.Directory,
		options.Description,
		options.Org,
		strings.Join(options.Platforms, ","),
		overwriteOption,
		pubGetOption,
		options.IosLang,
		options.AndroidLang,
	}

	itemOptionsBlock := lipgloss.JoinVertical(lipgloss.Left, itemOptions...)

	items := strings.Split(
		lipgloss.JoinHorizontal(lipgloss.Top, itemLabelsBlock, itemOptionsBlock),
		"\n",
	)

	choose, err := choose.New(choose.Options{
		Title: "Options",
		Items: items,
		Selected: []int{options.ChooseSelected},
		Limit: true,
		DisableHelp: true,
	})
	if err != nil {
		return Model{}, err
	}

	choosebutton := choosebutton.New(choosebutton.Options{
		Choice: options.ConfirmChoice,
		Prompt: "Confirm Options",
		Items: []string{"Yes", "No"},
		DisableHelp: true,
		Center: true,
	})

	return Model{
		choose: choose,
		choosebutton: choosebutton,
		help: help.New(),
		keyMap: keyMap{
			optionsSelected: options.OptionsSelected,
			additionalShortHelpKeys: func(optionsSelected bool) []key.Binding {
				if optionsSelected {
					return choose.List.ShortHelp()
				}
				return choosebutton.KeyMap.ShortHelp()
			},
			switchSelected: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "switch"),
			),
		},
		OptionsSelected: options.OptionsSelected,
		ChooseSelected: options.ChooseSelected,
		ConfirmChoice: options.ConfirmChoice,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.switchSelected) {
			m.OptionsSelected = !m.OptionsSelected
		}
	}
	m.keyMap.optionsSelected = m.OptionsSelected

	var cmd tea.Cmd
	if m.OptionsSelected {
		var model tea.Model
		model, cmd = m.choose.Update(msg)
		m.choose = model.(choose.Model)
	} else {
		m.choosebutton, cmd = m.choosebutton.Update(msg)
	}

	m.ChooseSelected = m.choose.List.Index()
	m.ConfirmChoice  = m.choosebutton.Choice

	return m, cmd
}

func (m Model) View() string {
	// remove extra line from options
	optionsSplit := strings.Split(m.choose.View(), "\n")
	options := strings.Join(optionsSplit[:len(optionsSplit)-1], "\n")
	
	confirm := m.choosebutton.View()
	if m.OptionsSelected {
		options = selectedStyle.Render(options)
		confirm = unselectedStyle.Render(confirm)
	} else {
		options = unselectedStyle.Render(options)
		confirm = selectedStyle.Render(confirm)
	}

	return constants.ViewStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		constants.PromptStyle.Render("Confirm Flutter Project Options"),
		lipgloss.JoinHorizontal(lipgloss.Top, confirm, options),
		m.help.View(m.keyMap),
	))
}