package selector

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

type state int

const STATE_PENDING state = 1
const STATE_COMPLETE state = 2
const STATE_ERROR state = 3

func New(width, height int) Model {
	a := Model{}
	a.table = table.New(table.WithWidth(width), table.WithHeight(height-2), table.WithFocused(true))
	a.spinner = spinner.New()
	a.spinner.Spinner = spinner.Points
	a.state = STATE_PENDING
	return a
}

type Model struct {
	networks []network
	error    string
	table    table.Model
	selected string
	spinner  spinner.Model
	state    state
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(run, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "<ctrl>+c":
			return m, tea.Quit
		case "enter":
			m.selected = m.table.SelectedRow()[1]
		}
	case nmcliSuccess:
		rows := make([]table.Row, len(msg.networks))
		for i, network := range msg.networks {
			rows[i] = table.Row{network.inUse, network.ssid, network.signal}
		}
		cols := []table.Column{{Title: "In Use", Width: 6}, {Title: "SSID", Width: 50}, {Title: "Signal", Width: 6}}
		m.table.SetColumns(cols)
		m.table.SetRows(rows)
		m.state = STATE_COMPLETE
	case nmcliFailed:
		m.error = msg.err.Error()
		m.state = STATE_ERROR
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	m.spinner, cmd = m.spinner.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, cmd
}

func (m Model) View() string {
	if m.state == STATE_PENDING {
		return m.spinner.View()
	}
	if m.state == STATE_ERROR {
		return m.error
	}
	if m.state == STATE_COMPLETE {
		return baseStyle.Render(m.table.View())
	}
	panic("unknown state")
}

func (m Model) Selected() string {
	return m.table.SelectedRow()[1]
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
