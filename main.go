package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const basic = `This is a basic test text`

type model struct {
	text        string
	placeholder string
	cursor      int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, tea.Quit
		case "backspace":
			if m.cursor != 0 {
				m.cursor--
				m.text = m.text[:len(m.text)-1]
			}
		default:
			m.cursor++
			m.text += msg.String()
		}
	}

	return m, nil
}

func (m model) View() string {
	return m.text + m.placeholder[m.cursor:]
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		text:        "",
		placeholder: basic,
		cursor:      0,
	}
}
