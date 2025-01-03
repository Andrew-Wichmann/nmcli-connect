package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const STATE_SELECTING = 1
const STATE_PROMPT = 2
const STATE_CONNECT = 3

type app struct {
	state         state
	ssid          string
	selector      selector
	passwordInput passwordInput
	connector     connector
}

func newApp() app {
	a := app{selector: newSelector(), state: STATE_SELECTING}
	return a
}

func (a app) Init() tea.Cmd {
	return a.selector.Init()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "enter":
			if a.state == STATE_SELECTING {
				a.state = STATE_PROMPT
				a.ssid = a.selector.table.SelectedRow()[1]
				a.passwordInput = newPasswordInput(a.ssid)
				return a, a.passwordInput.Init()
			} else if a.state == STATE_PROMPT {
				a.state = STATE_CONNECT
				password := a.passwordInput.ti.Value()
				a.connector = newConnector(a.ssid, password)
				return a, a.connector.Init()
			}
		}
	}
	var cmd tea.Cmd
	if a.state == STATE_SELECTING {
		a.selector, cmd = a.selector.Update(msg)
	} else if a.state == STATE_PROMPT {
		a.passwordInput, cmd = a.passwordInput.Update(msg)
	} else if a.state == STATE_CONNECT {
		a.connector, cmd = a.connector.Update(msg)
	}
	return a, cmd
}

func (a app) View() string {
	if a.state == STATE_SELECTING {
		return a.selector.View()
	} else if a.state == STATE_PROMPT {
		return a.passwordInput.View()
	} else if a.state == STATE_CONNECT {
		return a.connector.View()
	}
	panic(fmt.Errorf("Unknown state"))
}
