package passwordinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textinput textinput.Model
}

func New(ssid string) Model {
	m := Model{}
	m.textinput = textinput.New()
	m.textinput.Prompt = fmt.Sprintf("Password for %s: ", ssid)
	m.textinput.EchoMode = textinput.EchoPassword
	m.textinput.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return m.textinput.Cursor.BlinkCmd()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.textinput.View()
}

func (m Model) Password() string {
	return m.textinput.Value()
}
