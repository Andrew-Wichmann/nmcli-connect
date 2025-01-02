package main

import tea "github.com/charmbracelet/bubbletea"

type app struct{}

func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "<ctrl>+c":
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a app) View() string {
	return "Hello World"
}
