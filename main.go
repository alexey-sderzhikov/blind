package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const basic = `This is a basic test text`

var (
	mistakeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Background(lipgloss.Color("196"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)

type model struct {
	text        string
	placeholder string
	cursor      int
	mistakes    map[int]bool
	currentView string
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
				if m.mistakes[m.cursor] {
					delete(m.mistakes, m.cursor)
				}
				m.text = m.text[:len(m.text)-1]
			}
		default:
			// TODO should use runes
			if string(m.placeholder[m.cursor]) != msg.String() {
				m.mistakes[m.cursor] = true
			}
			m.cursor++
			m.text += msg.String()
		}
	}

	return m, nil
}

func (m model) View() string {
	var view string
	for i, r := range m.text {
		if m.mistakes[i] {
			view += mistakeStyle.Render(string(m.placeholder[i]))
		} else {
			view += textStyle.Render(string(r))
		}
	}

	return view + placeholderStyle.Render(m.placeholder[m.cursor:])
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
		mistakes:    make(map[int]bool),
	}
}
