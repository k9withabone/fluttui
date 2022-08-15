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

type Config struct {
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
	PageSelected    int
	ConfirmChoice   int
}

type option struct {
	label    string
	value    string
	page     int
	position int
}

type Model struct {
	options      []option
	choose       choose.Model
	choosebutton choosebutton.Model
	help         help.Model
	keyMap       keyMap

	OptionsSelected bool
	PageSelected    int
	ConfirmChoice   int
}

func New(config Config) (Model, error) {
	options := []option{
		{
			label: "Template: ",
			value: config.Template,
			page: constants.Pages["templateSelect"],
		},
		{
			label: "Name: ",
			value: config.Name,
			page: constants.Pages["nameInput"],
		},
		{
			label: "Directory: ",
			value: config.Directory,
			page: constants.Pages["directoryInput"],
		},
		{
			label: "Description: ",
			value: config.Description,
			page: constants.Pages["descriptionInput"],
		},
		{
			label: "Orginization: ",
			value: config.Org,
			page: constants.Pages["orgInput"],
		},
	}

	if len(config.Platforms) != 0 {
		options = append(options, option{
			label: "Platforms: ",
			value: strings.Join(config.Platforms, ","),
			page: constants.Pages["platformSelect"],
		})
	}

	overwriteOption := "No"
	if config.Overwrite {
		overwriteOption = "Yes"
	}

	pubGetOption := "Yes"
	switch config.PubGet {
	case pubget.PubGetOffline:
		pubGetOption = "Offline"
	case pubget.PubGetNo:
		pubGetOption = "No"
	}

	options = append(
		options,
		option{
			label: "Overwrite: ",
			value: overwriteOption,
			page: constants.Pages["overwriteSelect"],
		},
		option{
			label: "pub get: ",
			value: pubGetOption,
			page: constants.Pages["pubGetSelect"],
		},
	)

	var platformIos, platformAndroid bool
	for _, platform := range config.Platforms {
		if platform == "ios" {
			platformIos = true
		} else if platform == "android" {
			platformAndroid = true
		}
	}

	if platformIos {
		options = append(options, option{
			label: "iOS Language: ",
			value: config.IosLang,
			page: constants.Pages["iosLangSelect"],
		})
	}

	if platformAndroid {
		options = append(options, option{
			label: "Android Language: ",
			value: config.AndroidLang,
			page: constants.Pages["androidLangSelect"],
		})
	}

	var labels []string
	var values []string
	var selected int
	for i, option := range options {
		labels = append(labels, option.label)
		values = append(values, option.value)
		if option.page == config.PageSelected {
			selected = i
		}
	}

	labelsBlock := lipgloss.JoinVertical(lipgloss.Left, labels...)
	valuesBlock := lipgloss.JoinVertical(lipgloss.Left, values...)

	items := strings.Split(
		lipgloss.JoinHorizontal(lipgloss.Top, labelsBlock, valuesBlock),
		"\n",
	)

	choose, err := choose.New(choose.Options{
		Title: "Options",
		Items: items,
		Selected: []int{selected},
		Limit: true,
		DisableHelp: true,
	})
	if err != nil {
		return Model{}, err
	}

	choosebutton := choosebutton.New(choosebutton.Options{
		Choice: config.ConfirmChoice,
		Prompt: "Confirm Options",
		Items: []string{"Yes", "No"},
		DisableHelp: true,
		Center: true,
	})

	return Model{
		options: options,
		choose: choose,
		choosebutton: choosebutton,
		help: help.New(),
		keyMap: keyMap{
			optionsSelected: config.OptionsSelected,
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
		OptionsSelected: config.OptionsSelected,
		PageSelected: config.PageSelected,
		ConfirmChoice: config.ConfirmChoice,
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

	m.PageSelected  = m.options[m.choose.List.Index()].page
	m.ConfirmChoice = m.choosebutton.Choice

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