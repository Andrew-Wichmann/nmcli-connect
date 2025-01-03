package main

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func newConnector(ssid, password string) connector {
	s := spinner.New()
	s.Spinner = spinner.Points
	return connector{ssid: ssid, password: password, spinner: s}
}

type connector struct {
	ssid     string
	password string
	error    string
	message  string
	spinner  spinner.Model
}

func (c connector) Init() tea.Cmd {
	return tea.Batch(c.connect, c.spinner.Tick)
}

func (c connector) Update(msg tea.Msg) (connector, tea.Cmd) {
	if msg, ok := msg.(connectFailed); ok {
		c.error = msg.err.Error()
	}
	if _, ok := msg.(connectSucceeded); ok {
		c.message = "Connected!"
	}
	var cmd tea.Cmd
	c.spinner, cmd = c.spinner.Update(msg)
	return c, cmd
}

func (c connector) View() string {
	if c.error != "" {
		return c.error
	}
	if c.message != "" {
		return c.message
	}
	return fmt.Sprintf("%s - connecting", c.spinner.View())
}

type connectFailed struct {
	err error
}

type connectSucceeded struct{}

func (c connector) connect() tea.Msg {
	cmd := exec.Command("nmcli", "device", "wifi", "connect", c.ssid, "password", c.password)
	_, err := cmd.Output()
	if err != nil {
		return connectFailed{err: err}
	}
	return connectSucceeded{}
}
