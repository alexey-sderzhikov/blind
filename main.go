package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const basic = `This is a basic test text for training with blind typing! This is a basic test text for training with blind typing! This is a basic test text for training with blind typing! This is a basic test text for training with blind typing! This is a basic test text for training with blind typing! This is a basic test text for training with blind typing! This is a basic test text for training with blind typing!This is a basic test text for training with blind typing!`

var (
	mistakeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("234")).Background(lipgloss.Color("202"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Bold(true)
	cursorStyle      = lipgloss.NewStyle().Background(lipgloss.Color("248"))
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	borderStyle      = lipgloss.NewStyle().Width(100)
)

type model struct {
	text        []rune
	placeholder []rune
	cursor      int
	mistakes    map[int]bool
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
			if len(msg.Runes) == 1 {
				if m.placeholder[m.cursor] != msg.Runes[0] {
					m.mistakes[m.cursor] = true
				}
				m.cursor++
				m.text = append(m.text, msg.Runes...)
			}
		}
	}

	return m, nil
}

func calcPercentageAcc(text, mistakes int) int {
	tperc := float32(text) / 100
	mperc := float32(mistakes) / tperc
	return int(100 - mperc)
}

func (m model) View() string {
	title := titleStyle.Render(fmt.Sprintf("Mistakes: %d %d%%", len(m.mistakes), calcPercentageAcc(len(m.text), len(m.mistakes))))
	var view = title + "\n"
	for i, r := range m.text {
		if m.mistakes[i] {
			view += mistakeStyle.Render(string(m.placeholder[i]))
		} else {
			view += textStyle.Render(string(r))
		}
	}

	view += cursorStyle.Render(string(m.placeholder[m.cursor])) + placeholderStyle.Render(string(m.placeholder[m.cursor+1:]))

	return borderStyle.Render(view)
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
		text:        make([]rune, 0),
		placeholder: []rune(basic),
		cursor:      0,
		mistakes:    make(map[int]bool),
	}
}
