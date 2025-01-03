package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type passwordInput struct {
	ti textinput.Model
}

func newPasswordInput(ssid string) passwordInput {
	pi := passwordInput{}
	pi.ti = textinput.New()
	pi.ti.Prompt = fmt.Sprintf("Password for %s: ", ssid)
	pi.ti.EchoMode = textinput.EchoPassword
	pi.ti.Focus()
	return pi
}

func (pi passwordInput) Init() tea.Cmd {
	return pi.ti.Cursor.BlinkCmd()
}

func (pi passwordInput) Update(msg tea.Msg) (passwordInput, tea.Cmd) {
	var cmd tea.Cmd
	pi.ti, cmd = pi.ti.Update(msg)
	return pi, cmd
}

func (pi passwordInput) View() string {
	return pi.ti.View()
}
