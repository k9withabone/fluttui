package choose

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/k9withabone/fluttui/tui/constants"
)

var (
	selectedColor     = constants.FlutterSky
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	hoverItemStyle    = lipgloss.NewStyle().PaddingLeft(2).Foreground(selectedColor)
	selectedItemStyle = itemStyle.Copy().Foreground(selectedColor)
)

type keyMap struct {
	Select key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Select,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.Select,
	}}
}

type item struct{
	text     string
	selected bool
}

func (i item) FilterValue() string {
	return i.text
}

type singleSelectDelegate struct{}
type multiSelectDelegate  struct{}

func delegateHeight(d list.ItemDelegate) int { return 1 }
func (d singleSelectDelegate) Height()   int { return delegateHeight(d) }
func (d multiSelectDelegate)  Height()   int { return delegateHeight(d) }

func delegateSpacing(d list.ItemDelegate) int { return 0 }
func (d singleSelectDelegate) Spacing()   int { return delegateSpacing(d) }
func (d multiSelectDelegate)  Spacing()   int { return delegateSpacing(d) }

func delegateUpdate(d list.ItemDelegate, msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d singleSelectDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return delegateUpdate(d, msg, m)
}
func (d multiSelectDelegate)  Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return delegateUpdate(d, msg, m)
}

func (d singleSelectDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.text)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return hoverItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}
func (d multiSelectDelegate)  Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := "[ ]"
	if i.selected && index == m.Index() {
		str = "> [x]"
	} else if i.selected {
		str = "[x]"
	} else if index == m.Index() {
		str = "> [•]"
	}

	if index == m.Index() {
		str = fmt.Sprintf("%s %s", str, i.text)
		str = hoverItemStyle.Render(str)
	} else if i.selected {
		str = selectedItemStyle.Render(str)
		str = fmt.Sprintf("%s %s", str, i.text)
	} else {
		str = fmt.Sprintf("%s %s", str, i.text)
		str = itemStyle.Render(str)
	}

	fmt.Fprintf(w, str)
}

type Options struct {
	Title       string
	Items       []string
	Selected    []int
	Limit       bool
	DisableHelp bool
}

type Model struct {
	Choices []string
	limit   bool
	List    list.Model
	keyMap  keyMap
}

func New(options Options) (Model, error) {
	itemsLen := len(options.Items)
	for _, index := range options.Selected {
		if index < 0 || index >= itemsLen{
			return Model{}, errors.New("a selected index is outside range of items")
		}
	}

	var listItems []list.Item

	for i, listItem := range options.Items {
		selected := false
		if !options.Limit {
			for _, index := range options.Selected {
				if i == index {
					selected = true
				}
			}
		}
		listItems = append(listItems, item{
			text: listItem,
			selected: selected,
		})
	}

	var d list.ItemDelegate
	if options.Limit {
		d = singleSelectDelegate{}
	} else {
		d = multiSelectDelegate{}
	}

	l := list.New(listItems, d, 80, itemsLen + 4)
	l.Title = options.Title
	l.Styles.Title = constants.PromptStyle

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()
	l.SetShowHelp(!options.DisableHelp)
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.KeyMap.CloseFullHelp.SetEnabled(false)

	if len(options.Selected) > 0 {
		l.Select(options.Selected[0])
	}

	k := keyMap{}
	if !options.Limit {
		k.Select = key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		)
	}

	l.AdditionalShortHelpKeys = func() []key.Binding {
		if options.Limit {
			return constants.Keys.ShortHelp()
		}
		return append(k.ShortHelp(), constants.Keys.ShortHelp()...)
	}

	m := Model{
		limit: options.Limit,
		List: l,
		keyMap: k,
	}

	// update to populate m.Choices
	model, _ := m.Update(nil)
	newModel, ok := model.(Model)
	if !ok {
		return m, errors.New("Error creating choose")
	}
	m = newModel

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var selectCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.Help.Width = msg.Width
		return m, nil
	
	case tea.KeyMsg:
		if !m.limit && key.Matches(msg, m.keyMap.Select) {
			index := m.List.Index()
			item, ok := m.List.Items()[index].(item)
			if ok {
				item.selected = !item.selected
				selectCmd = m.List.SetItem(index, item)
			}
		}
	}

	var updateCmd tea.Cmd
	m.List, updateCmd = m.List.Update(msg)

	m.Choices = []string{}
	currentIndex := m.List.Index()
	for index, listItem := range m.List.Items() {
		listItem, ok := listItem.(item)
		if !ok {
			continue
		}
		if index == currentIndex || listItem.selected {
			// append "hovered" and selected items
			m.Choices = append(m.Choices, listItem.text)
		}
	}

	return m, tea.Batch(selectCmd, updateCmd)
}

func (m Model) View() string {
	return "\n" + m.List.View()
}