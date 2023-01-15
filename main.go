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
	"github.com/tjarratt/babble"
)

var (
	mistakeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("234")).Background(lipgloss.Color("202"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Bold(true)
	cursorStyle      = lipgloss.NewStyle().Background(lipgloss.Color("87"))
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	borderStyle      = lipgloss.NewStyle().Width(100)
)

type window int

const (
	menu window = iota
	typing
	results
	mode
)

func (w window) String() string {
	switch w {
	case 0:
		return "menu"
	case 1:
		return "typing"
	case 2:
		return "results"
	case 3:
		return "mode"
	}
	return "unknown"
}

type model struct {
	currendWindow window

	texts []string // loaded texts

	text        []rune // user's typed text
	placeholder []rune // expected text

	cursor int // define current symbol in placeholder

	mistakes     map[int]bool // all indexes of mistakes which did user during typing,
	mistakeCount int

	typedSybmolsCount int
	typedStartTime    time.Time
	typedEndTime      time.Time
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
			m.typedEndTime = time.Now()
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
			if len(msg.Runes) != 1 {
				break
			}

			m.typedSybmolsCount++
			// start timer when typed first symbol
			if m.typedSybmolsCount == 1 {
				m.typedStartTime = time.Now()
			}

			// detect mistake if typed sybmol != expected symbol
			if m.placeholder[m.cursor] != msg.Runes[0] {
				m.mistakeCount++
				m.mistakes[m.cursor] = true
			}

			m.cursor++
			m.text = append(m.text, msg.Runes[0])

			// if typed last symbol then end the test
			if m.cursor == len(m.placeholder) {
				m.typedEndTime = time.Now()
				m.currendWindow = results
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
			m.cursor = 0
			m.currendWindow = mode
		}
	}

	return m, nil
}

func (m model) updateModeWindow(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "up":
			if m.cursor == 1 {
				m.cursor = 0
			}
		case "down":
			if m.cursor == 0 {
				m.cursor = 1
			}
		case "enter":
			if m.cursor == 0 {
				m.placeholder = generateRandomWords()
			} else {
				var err error
				m.placeholder, err = chooseRandomText(m.texts)
				if err != nil {
					panic(err)
				}
			}

			m.pruneForNewText()
			m.currendWindow = typing
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currendWindow {
	case typing:
		return m.updateTypingWindow(msg)
	case results:
		return m.updateResultsWindow(msg)
	case mode:
		return m.updateModeWindow(msg)
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
	case mode:
		return m.viewMode()
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
	t := m.typedEndTime.Sub(m.typedStartTime).Round(time.Second)
	return fmt.Sprintf("Mistakes: %d\nTyped symbols: %d\nAccuracy: %d%%\nTime: %s",
		m.mistakeCount,
		m.typedSybmolsCount,
		calcPercentageAcc(m.typedSybmolsCount, m.mistakeCount),
		t.String(),
	)
}

func (m model) viewMode() string {
	var view string
	if m.cursor == 0 {
		view = "> words\ntext"
	} else {
		view = "words\n> text"
	}

	return view
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
		currendWindow: mode,
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

func generateRandomWords() []rune {
	bab := babble.NewBabbler()
	bab.Count = 10
	bab.Separator = " "
	return []rune(bab.Babble())
}
