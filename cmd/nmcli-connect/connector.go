package main

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func newConnector(ssid, password string) connector {
	return connector{ssid: ssid, password: password}
}

type connector struct {
	ssid     string
	password string
	error    string
	message  string
}

func (c connector) Init() tea.Cmd {
	return c.connect
}

func (c connector) Update(msg tea.Msg) (connector, tea.Cmd) {
	if msg, ok := msg.(connectFailed); ok {
		c.error = msg.err.Error()
	}
	if _, ok := msg.(connectSucceeded); ok {
		c.message = "Connected!"
	}
	return c, nil
}

func (c connector) View() string {
	if c.error != "" {
		return c.error
	}
	if c.message != "" {
		return c.message
	}
	return fmt.Sprintf("beep boop. connecting to %s with password: %s", c.ssid, c.password)
}

type connectFailed struct {
	err error
}

type connectSucceeded struct{}

func (c connector) connect() tea.Msg {
	cmd := exec.Command("sleep", "10")
	_, err := cmd.Output()
	if err != nil {
		return connectFailed{err: err}
	}
	return connectSucceeded{}
}
