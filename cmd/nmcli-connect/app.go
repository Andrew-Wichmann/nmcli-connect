package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	output string
}

type start struct{}

func (a app) Init() tea.Cmd {
	return func() tea.Msg { return start{} }
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "<ctrl>+c":
			return a, tea.Quit
		}
	case start:
		cmd := exec.Command("nmcli", "device", "wifi")
		output, err := cmd.Output()
		if err != nil {
			a.output = string(err.Error())
		} else {
			a.output = string(output)
		}
	}
	return a, nil
}

func (a app) View() string {
	return a.output
}
