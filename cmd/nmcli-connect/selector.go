package main

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectorState int

const STATE_PENDING selectorState = 1
const STATE_COMPLETE selectorState = 2
const STATE_ERROR selectorState = 3

func newSelector() selector {
	a := selector{}
	a.table = table.New(table.WithWidth(100), table.WithHeight(50), table.WithFocused(true))
	a.spinner = spinner.New()
	a.spinner.Spinner = spinner.Points
	a.state = STATE_PENDING
	return a
}

type selector struct {
	networks []network
	error    string
	table    table.Model
	selected string
	spinner  spinner.Model
	state    selectorState
}

func (a selector) Init() tea.Cmd {
	return tea.Batch(run, a.spinner.Tick)
}

func (a selector) Update(msg tea.Msg) (selector, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "<ctrl>+c":
			return a, tea.Quit
		case "enter":
			a.selected = a.table.SelectedRow()[1]
		}
	case nmcliSuccess:
		rows := make([]table.Row, len(msg.networks))
		for i, network := range msg.networks {
			rows[i] = table.Row{network.inUse, network.ssid, network.signal}
		}
		cols := []table.Column{{Title: "In Use", Width: 5}, {Title: "SSID", Width: 50}, {Title: "Signal", Width: 10}}
		a.table.SetColumns(cols)
		a.table.SetRows(rows)
		a.state = STATE_COMPLETE
	case nmcliFailed:
		a.error = msg.err.Error()
		a.state = STATE_ERROR
	}
	var cmd tea.Cmd
	a.table, cmd = a.table.Update(msg)
	if cmd != nil {
		return a, cmd
	}
	a.spinner, cmd = a.spinner.Update(msg)
	if cmd != nil {
		return a, cmd
	}
	return a, cmd
}

func (a selector) View() string {
	if a.state == STATE_PENDING {
		return a.spinner.View()
	}
	if a.state == STATE_ERROR {
		return a.error
	}
	if a.state == STATE_COMPLETE {
		return baseStyle.Render(a.table.View())
	}
	panic("unknown state")
}

type network struct {
	inUse  string
	ssid   string
	signal string
}

type nmcliFailed struct {
	err error
}
type nmcliSuccess struct {
	networks []network
}

func run() tea.Msg {
	cmd := exec.Command("nmcli", "device", "wifi")
	output, err := cmd.Output()
	if err != nil {
		return nmcliFailed{err: err}
	}
	networks, err := parseOutput(output)
	if err != nil {
		return nmcliFailed{err: err}
	}
	return nmcliSuccess{networks: networks}
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
