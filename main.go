package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
