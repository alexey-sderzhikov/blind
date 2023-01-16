package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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
