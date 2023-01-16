package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tjarratt/babble"
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

func (m model) loadTexts() ([]string, error) {
	b, err := os.ReadFile("texts")
	if err != nil {
		return nil, fmt.Errorf("error during reading text from file: %v", err)
	}

	return strings.Split(strings.Trim(string(b), "\n"), "\n"), nil
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
