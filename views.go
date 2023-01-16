package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	mistakeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("234")).Background(lipgloss.Color("202"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Bold(true)
	cursorStyle      = lipgloss.NewStyle().Background(lipgloss.Color("87"))
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	borderStyle      = lipgloss.NewStyle().Width(100)
)

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

func calcPercentageAcc(text, mistakes int) int {
	if text == 0 {
		return 0
	}

	tperc := float32(text) / 100
	mperc := float32(mistakes) / tperc
	return int(100 - mperc)
}
