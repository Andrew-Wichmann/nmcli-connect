package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	prog := tea.NewProgram(app{})
	_, err := prog.Run()
	if err != nil {
		panic(err)
	}
}
