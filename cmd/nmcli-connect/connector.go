package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func newConnector(ssid, password string) connector {
	return connector{ssid: ssid, password: password}
}

type connector struct {
	ssid     string
	password string
}

func (c connector) Init() tea.Cmd {
	return nil
}

func (c connector) Update(msg tea.Msg) (connector, tea.Cmd) {
	return c, nil
}

func (c connector) View() string {
	return fmt.Sprintf("beep boop. connecting to %s with password: %s", c.ssid, c.password)
}
