package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	networks []network
	error    string
}

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
			a.error = string(err.Error())
		} else {
			networks, err := parseOutput(output)
			if err != nil {
				a.error = err.Error()
			} else {
				a.networks = networks
			}
		}
	}
	return a, nil
}

func (a app) View() string {
	if a.error != "" {
		return a.error
	}
	var builder strings.Builder
	for _, network := range a.networks {
		builder.WriteString(fmt.Sprintf("isUse: %t, SSID: %s, signal: %d\n", network.inUse, network.ssid, network.signal))
	}
	return builder.String()
}

type start struct{}

type network struct {
	inUse  bool
	ssid   string
	signal int
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
	for _, line := range lines[1 : len(lines)-1] {
		fields := make([]string, 8)
		for i, col := range cols {
			fields[i] = strings.Trim(line[col.start:col.end], " ")
		}
		var inUse bool
		if strings.Trim(fields[0], " ") == "*" {
			inUse = true
		}
		ssid := strings.Trim(fields[2], " ")
		signal, err := strconv.Atoi(fields[6])
		if err != nil {
			return nil, fmt.Errorf("Could not parse signal")
		}
		networks = append(networks, network{inUse: inUse, ssid: ssid, signal: signal})
	}
	return networks, nil
}
