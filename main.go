package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	mistakeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("234")).Background(lipgloss.Color("202"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Bold(true)
	cursorStyle      = lipgloss.NewStyle().Background(lipgloss.Color("99"))
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	borderStyle      = lipgloss.NewStyle().Width(100)
)

type window int

const (
	menu window = iota
	typing
	results
)

func (w window) String() string {
	switch w {
	case 0:
		return "menu"
	case 1:
		return "typing"
	case 2:
		return "results"
	}
	return "unknown"
}

type model struct {
	currendWindow     window
	text              []rune
	placeholder       []rune
	cursor            int
	mistakes          map[int]bool
	mistakeCount      int
	typedSybmolsCount int
	texts             []string
}

func (m *model) pruneForNewText() {
	m.cursor = 0
	m.mistakeCount = 0
	m.mistakes = make(map[int]bool)
	m.typedSybmolsCount = 0
	m.text = make([]rune, 0)
}

func (m model) Init() tea.Cmd {
	return nil
}

// navigation logic base for most pages
func (m model) navigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		return nil, tea.Quit
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		// if m.cursor < m.objectCount-1 {
		// 	m.cursor++
		// }
	}

	return m, nil
}

func (m model) updateMenuWindow(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return nil, tea.Quit
		default:
			return m.navigation(msg)
		}
	}

	return m, nil
}

func (m model) updateTypingWindow(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "tab":
			var err error
			m.placeholder, err = chooseRandomText(m.texts)
			if err != nil {
				panic(err)
			}

			m.pruneForNewText()
		case "enter":
			m.currendWindow = results
		case "backspace":
			if m.cursor != 0 {
				m.cursor--
				if m.mistakes[m.cursor] {
					delete(m.mistakes, m.cursor)
				}
				m.text = m.text[:len(m.text)-1]
			}
		default:
			m.typedSybmolsCount++
			if len(msg.Runes) == 1 {
				if m.placeholder[m.cursor] != msg.Runes[0] {
					m.mistakeCount++
					m.mistakes[m.cursor] = true
				}
				m.cursor++
				m.text = append(m.text, msg.Runes[0])

				if m.cursor == len(m.placeholder)-1 {
					m.currendWindow = results
				}
			}
		}
	}

	return m, nil
}

func (m model) updateResultsWindow(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "tab", "enter":
			var err error
			m.placeholder, err = chooseRandomText(m.texts)
			if err != nil {
				panic(err)
			}
			m.pruneForNewText()
			m.currendWindow = typing
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currendWindow {
	case menu:
		return m.updateMenuWindow(msg)
	case typing:
		return m.updateTypingWindow(msg)
	case results:
		return m.updateResultsWindow(msg)
	}

	return m, nil
}

func calcPercentageAcc(text, mistakes int) int {
	if text == 0 {
		return 0
	}

	tperc := float32(text) / 100
	mperc := float32(mistakes) / tperc
	return int(100 - mperc)
}

func (m model) View() string {
	switch m.currendWindow {
	case typing:
		return m.viewTyping()
	case results:
		return m.viewResults()
	}
	return fmt.Sprintf("not found suitable window for [%s]", m.currendWindow.String())
}

func (m model) viewTyping() string {
	title := titleStyle.Render(
		fmt.Sprintf("Mistakes: %d %d%%",
			m.mistakeCount,
			calcPercentageAcc(m.typedSybmolsCount, m.mistakeCount),
		),
	)

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

func (m model) viewResults() string {
	return fmt.Sprintf("Mistakes: %d\nTyped symbols: %d\nAccuracy: %d%%",
		m.mistakeCount, m.typedSybmolsCount, calcPercentageAcc(m.typedSybmolsCount, m.mistakeCount))
}

func (m model) loadTexts() ([]string, error) {
	b, err := os.ReadFile("texts")
	if err != nil {
		return nil, fmt.Errorf("error during reading text from file: %v", err)
	}

	return strings.Split(strings.Trim(string(b), "\n"), "\n"), nil
}

func main() {
	p := tea.NewProgram(initialModel())
	p.Run()
}

func initialModel() model {
	m := model{
		currendWindow: typing,
		text:          make([]rune, 0),
		cursor:        0,
		mistakes:      make(map[int]bool),
		mistakeCount:  0,
	}

	var err error
	m.texts, err = m.loadTexts()
	if err != nil {
		panic(err)
	}

	m.placeholder, err = chooseRandomText(m.texts)
	if err != nil {
		panic(err)
	}

	return m
}

func chooseRandomText(texts []string) ([]rune, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	text := []rune(texts[r1.Intn(len(texts))])
	if len(text) == 0 {
		return nil, errors.New("chosen text is empty")
	}

	return text, nil
}
