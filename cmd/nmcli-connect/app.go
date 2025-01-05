package main

import (
	"fmt"

	"github.com/Andrew-Wichmann/nmcli-connect/internal/connector"
	"github.com/Andrew-Wichmann/nmcli-connect/internal/passwordinput"
	"github.com/Andrew-Wichmann/nmcli-connect/internal/selector"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const STATE_INIT = 1
const STATE_SELECTING = 2
const STATE_PROMPT = 3
const STATE_CONNECT = 4

type app struct {
	state         state
	ssid          string
	selector      selector.Model
	passwordInput passwordinput.Model
	connector     connector.Model
}

func newApp() app {
	a := app{state: STATE_INIT}
	return a
}

func (a app) Init() tea.Cmd {
	return a.selector.Init()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.selector = selector.New(msg.Width, msg.Height)
		a.state = STATE_SELECTING
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "enter":
			if a.state == STATE_SELECTING {
				a.state = STATE_PROMPT
				a.ssid = a.selector.Selected()
				a.passwordInput = passwordinput.New(a.ssid)
				return a, a.passwordInput.Init()
			} else if a.state == STATE_PROMPT {
				a.state = STATE_CONNECT
				password := a.passwordInput.Password()
				a.connector = connector.New(a.ssid, password)
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
	if a.state == STATE_INIT {
		return "Starting..."
	} else if a.state == STATE_SELECTING {
		return a.selector.View()
	} else if a.state == STATE_PROMPT {
		return a.passwordInput.View()
	} else if a.state == STATE_CONNECT {
		return a.connector.View()
	}
	panic(fmt.Errorf("Unknown state"))
}
