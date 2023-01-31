package main

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
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
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}
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
		case "s":
			m.currendWindow = settings
		}

	}

	return m, nil
}

func (m model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		// Set focus to next input
		case "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currendWindow {
	case typing:
		return m.updateTypingWindow(msg)
	case results:
		return m.updateResultsWindow(msg)
	case mode:
		return m.updateModeWindow(msg)
	case settings:
		return m.updateSettings(msg)
	}

	return m, nil
}
