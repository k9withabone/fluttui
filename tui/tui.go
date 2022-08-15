package tui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k9withabone/fluttui/tui/components/choose"
	"github.com/k9withabone/fluttui/tui/components/input"
	"github.com/k9withabone/fluttui/tui/constants"
	"github.com/k9withabone/fluttui/tui/pages/androidlang"
	"github.com/k9withabone/fluttui/tui/pages/confirm"
	"github.com/k9withabone/fluttui/tui/pages/description"
	"github.com/k9withabone/fluttui/tui/pages/directory"
	"github.com/k9withabone/fluttui/tui/pages/ioslang"
	"github.com/k9withabone/fluttui/tui/pages/name"
	"github.com/k9withabone/fluttui/tui/pages/org"
	"github.com/k9withabone/fluttui/tui/pages/overwrite"
	"github.com/k9withabone/fluttui/tui/pages/platforms"
	"github.com/k9withabone/fluttui/tui/pages/pubget"
	"github.com/k9withabone/fluttui/tui/pages/run"
	"github.com/k9withabone/fluttui/tui/pages/template"
)

func StartTea() {
	m := mainModel{
		pubGet: pubget.PubGetYes,
		iosLang: "swift",
		androidLang: "kotlin",
		confirmVisited: false,
	}
	m.setPage(constants.Pages["templateSelect"])
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("There was an error: %v\n", err)
		os.Exit(1)
	}
}

type confrimOptions struct {
	optionsSelected bool
	pageSelected    int
	confirmChoice   int
}

type mainModel struct {
	page           int
	currentModel   tea.Model
	confirmVisited bool
	confrimOptions confrimOptions
	cmdString      string

	template       string
	name           string
	directory      string
	description    string
	org            string
	platforms      []string
	overwrite      bool
	pubGet         int
	iosLang        string
	androidLang    string
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, runCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, constants.Keys.Quit) ||
		   (!m.confirmVisited && m.page == 0 &&
			key.Matches(msg, constants.Keys.Back)) {
			// quit if quit key is pressed
			// or if at first page and back is pressed
			// and confirm has not been visited
			m.currentModel = nil
			return m, tea.Quit
		}
		if m.page < constants.Pages["running"] {
			var err error
			switch {
			case m.confirmVisited && m.page != constants.Pages["confirmation"] &&
			     (key.Matches(msg, constants.Keys.Back) ||
				  key.Matches(msg, constants.Keys.Next)):
				// go back to confirmation page if it has already
				// been visited and quit was not pressed
				cmd, err = m.setPage(constants.Pages["confirmation"])
	
			case key.Matches(msg, constants.Keys.Back) &&
			     m.page <= constants.Pages["confirmation"]:
				// go to previous page if at confirmation or before
				cmd, err = m.setPage(m.page - 1)
	
			case key.Matches(msg, constants.Keys.Next) &&
			     m.page < constants.Pages["confirmation"]:
				// go to next page if before confirmation
				cmd, err = m.setPage(m.page + 1)

			case key.Matches(msg, constants.Keys.Next) &&
			     m.page == constants.Pages["confirmation"]:
				// on confirmation page "next" key selects options to edit
				// or confirms the options and starts flutter create
				model := m.currentModel.(confirm.Model)
				if model.OptionsSelected {
					cmd, err = m.setPage(model.PageSelected)
				} else {
					if model.ConfirmChoice == 1 {
						// no confirm, shift focus to options
						model.OptionsSelected = true
						m.currentModel = model
					} else {
						// options confirmed
						cmd, m.cmdString = m.getFlutterCreateCmd()
						runCmd, _ = m.setPage(constants.Pages["running"])
					}
				}
			}
			if err != nil {
				println(err.Error())
				return m, tea.Quit
			}
		}

	case constants.CompletedMsg:
		if msg.Err != nil {
			println("There was an error running the command:\n", msg.Err.Error())
			return m, tea.Quit
		}
		cmd = tea.Quit
	}

	var newCmd tea.Cmd
	m.currentModel, newCmd = m.currentModel.Update(msg)

	return m, tea.Batch(cmd, runCmd, newCmd)
}

func (m mainModel) View() string {
	if m.currentModel != nil {
		return m.currentModel.View()
	}
	return ""
}

func (m *mainModel) setPage(page int) (tea.Cmd, error) {
	var err error
	var cmd tea.Cmd

	if page == constants.Pages["platformSelect"] &&
	   m.template != "app" &&
	   m.template != "plugin" &&
	   m.template != "plugin_ffi" {
		if page > m.page {
			page++
		} else {
			page--
		}
	}

	if m.currentModel != nil {
		// get entered info from the current page
		switch m.page {
		case constants.Pages["templateSelect"]:
			m.template = m.currentModel.(choose.Model).Choices[0]
			if m.template != "app" &&
			   m.template != "plugin" &&
			   m.template != "plugin_ffi" {
				m.platforms = []string{}
			} else if m.confirmVisited && len(m.platforms) == 0 {
				page = constants.Pages["platformSelect"]
			}
		
		case constants.Pages["nameInput"]:
			m.name = m.currentModel.(input.Model).Text
		
		case constants.Pages["directoryInput"]:
			m.directory = m.currentModel.(input.Model).Text
		
		case constants.Pages["descriptionInput"]:
			m.description = m.currentModel.(input.Model).Text
		
		case constants.Pages["orgInput"]:
			m.org = m.currentModel.(input.Model).Text
		
		case constants.Pages["platformSelect"]:
			m.platforms = m.currentModel.(choose.Model).Choices
		
		case constants.Pages["overwriteSelect"]:
			m.overwrite = m.currentModel.(overwrite.Model).Overwrite
		
		case constants.Pages["confirmation"]:
			model := m.currentModel.(confirm.Model)
			m.confrimOptions.optionsSelected = model.OptionsSelected
			m.confrimOptions.pageSelected = model.PageSelected
			m.confrimOptions.confirmChoice = model.ConfirmChoice
			constants.Keys.Next.SetHelp("enter", "confirm")
		
		case constants.Pages["pubGetSelect"]:
			m.pubGet = m.currentModel.(pubget.Model).Choice
		
		case constants.Pages["iosLangSelect"]:
			m.iosLang = m.currentModel.(ioslang.Model).Choice
		
		case constants.Pages["androidLangSelect"]:
			m.androidLang = m.currentModel.(androidlang.Model).Choice
		}
	}

	// set currentModel to apropriate model for new page
	switch page {
	case constants.Pages["templateSelect"]:
		m.currentModel, err = template.New(m.template)
	
	case constants.Pages["nameInput"]:
		m.currentModel, cmd, err = name.New(m.name)
	
	case constants.Pages["directoryInput"]:
		m.currentModel, cmd, err = directory.New(m.directory, m.name)
	
	case constants.Pages["descriptionInput"]:
		m.currentModel, cmd, err = description.New(m.description)
	
	case constants.Pages["orgInput"]:
		m.currentModel, cmd, err = org.New(m.org, m.name)
	
	case constants.Pages["platformSelect"]:
		m.currentModel, err = platforms.New(m.platforms)
	
	case constants.Pages["overwriteSelect"]:
		m.currentModel = overwrite.New(m.overwrite)
	
	case constants.Pages["confirmation"]:
		constants.Keys.Next.SetHelp("enter", "select")
		m.confirmVisited = true
		m.currentModel, err = confirm.New(confirm.Config{
			Template: m.template,
			Name: m.name,
			Directory: m.directory,
			Description: m.description,
			Org: m.org,
			Platforms: m.platforms,
			Overwrite: m.overwrite,
			PubGet: m.pubGet,
			IosLang: m.iosLang,
			AndroidLang: m.androidLang,

			OptionsSelected: m.confrimOptions.optionsSelected,
			PageSelected: m.confrimOptions.pageSelected,
			ConfirmChoice: m.confrimOptions.confirmChoice,
		})
	
	case constants.Pages["pubGetSelect"]:
		m.currentModel = pubget.New(m.pubGet)
	
	case constants.Pages["iosLangSelect"]:
		m.currentModel = ioslang.New(m.iosLang)
	
	case constants.Pages["androidLangSelect"]:
		m.currentModel = androidlang.New(m.androidLang)
	
	case constants.Pages["running"]:
		m.currentModel, cmd = run.New(m.cmdString)
	}

	m.page = page
	return cmd, err
}

func (m *mainModel) getFlutterCreateCmd() (tea.Cmd, string) {
	cmdArgs := []string{"create"}

	if m.directory == "" {
		cmdArgs = append(cmdArgs, m.name)
	} else {
		cmdArgs = append(cmdArgs, m.directory, "--project-name", m.name)
	}

	cmdArgs = append(cmdArgs, "-t", m.template)

	if len(m.platforms) != 0 {
		cmdArgs = append(cmdArgs, "--platforms", strings.Join(m.platforms, ","))
	}

	if m.description != "" {
		cmdArgs = append(cmdArgs, "--description", m.description)
	}

	if m.org != "" {
		cmdArgs = append(cmdArgs, "--org", m.org)
	}

	if m.overwrite {
		cmdArgs = append(cmdArgs, "--overwrite")
	}

	switch m.pubGet {
	case pubget.PubGetOffline:
		cmdArgs = append(cmdArgs, "--offline")
	case pubget.PubGetNo:
		cmdArgs = append(cmdArgs, "--no-pub")
	}

	var platformIos, platformAndroid bool
	for _, platform := range m.platforms {
		if platform == "ios" {
			platformIos = true
		} else if platform == "android" {
			platformAndroid = true
		}
	}

	if platformIos {
		cmdArgs = append(cmdArgs, "-i", m.iosLang)
	}

	if platformAndroid {
		cmdArgs = append(cmdArgs, "-a", m.androidLang)
	}

	cmd := exec.Command("flutter", cmdArgs...)

	cmdStr := strings.Join(cmd.Args, " ")

	return func() tea.Msg {
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		
		err := cmd.Run()

		return constants.CompletedMsg{
			Err: err,
			Output: out.String(),
		}
	}, cmdStr
}
