package main

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
	p.Run()
}

func initialModel() model {
	m := model{
		keys:          keys,
		help:          help.New(),
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
