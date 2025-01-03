package main

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newApp() app {
	a := app{}
	a.table = table.New(table.WithWidth(100), table.WithHeight(50), table.WithFocused(true))
	return a
}

type app struct {
	networks []network
	error    string
	table    table.Model
}

func (a app) Init() tea.Cmd {
	return func() tea.Msg { return start{} }
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.table, cmd = a.table.Update(msg)
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
			a.error = string(err.Error())
		} else {
			networks, err := parseOutput(output)
			if err != nil {
				a.error = err.Error()
			} else {
				a.networks = networks
			}
			rows := make([]table.Row, len(a.networks))
			for i, network := range a.networks {
				rows[i] = table.Row{network.inUse, network.ssid, network.signal}
			}
			cols := []table.Column{{Title: "In Use", Width: 5}, {Title: "SSID", Width: 50}, {Title: "Signal", Width: 10}}
			a.table.SetColumns(cols)
			a.table.SetRows(rows)
		}
	}
	return a, cmd
}

func (a app) View() string {
	if a.error != "" {
		return a.error
	}
	return baseStyle.Render(a.table.View())
}

type start struct{}

type network struct {
	inUse  string
	ssid   string
	signal string
}

func parseOutput(output []byte) ([]network, error) {
	var networks []network

	lines := strings.Split(string(output), "\n")

	if len(lines) < 2 {
		return nil, fmt.Errorf("No networks")
	}

	type collumn struct {
		start int
		end   int
	}
	cols := []collumn{}
	col := collumn{start: 0}
	prev := '-'
	for i, c := range lines[0] {
		if !unicode.IsSpace(c) && unicode.IsSpace(prev) {
			col.end = i - 1
			cols = append(cols, col)
			col.start = i
		}
		prev = c
	}
	for _, line := range lines[1 : len(lines)-1] { // Would be nice to not discard the last element
		fields := make([]string, 8)
		for i, col := range cols {
			fields[i] = strings.Trim(line[col.start:col.end], " ")
		}
		inUse := strings.Trim(fields[0], " ")
		ssid := strings.Trim(fields[2], " ")
		signal := strings.Trim(fields[6], " ")
		networks = append(networks, network{inUse: inUse, ssid: ssid, signal: signal})
	}
	return networks, nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("240"))
