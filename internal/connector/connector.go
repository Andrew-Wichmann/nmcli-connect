package connector

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func New(ssid, password string) Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	return Model{ssid: ssid, password: password, spinner: s}
}

type Model struct {
	ssid     string
	password string
	error    string
	message  string
	spinner  spinner.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.connect, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(connectFailed); ok {
		m.error = msg.err.Error()
	}
	if _, ok := msg.(connectSucceeded); ok {
		m.message = "Connected!"
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.error != "" {
		return m.error
	}
	if m.message != "" {
		return m.message
	}
	return fmt.Sprintf("%s - connecting", m.spinner.View())
}

type connectFailed struct {
	err error
}

type connectSucceeded struct{}

func (m Model) connect() tea.Msg {
	cmd := exec.Command("nmcli", "device", "wifi", "connect", m.ssid, "password", m.password)
	_, err := cmd.Output()
	if err != nil {
		return connectFailed{err: err}
	}
	return connectSucceeded{}
}
