package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tjarratt/babble"
)

type window int

const (
	menu window = iota
	typing
	results
	mode
	settings
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
	case 4:
		return "settings"
	}
	return "unknown"
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Chose    key.Binding
	Settings key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Settings, k.Quit, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Chose},      // first column
		{k.Help, k.Settings, k.Quit}, // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Chose: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "chose/finish"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
	Settings: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "settings"),
	),
}

type model struct {
	keys keyMap
	help help.Model

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

	focusIndex int
	inputs     map[string]textinput.Model
	cursorMode textinput.CursorMode
}

func (m *model) pruneForNewText() {
	m.cursor = 0
	m.mistakeCount = 0
	m.mistakes = make(map[int]bool)
	m.typedSybmolsCount = 0
	m.text = make([]rune, 0)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
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

func initialModel() model {
	m := model{
		keys:          keys,
		help:          help.New(),
		currendWindow: mode,
		text:          make([]rune, 0),
		cursor:        0,
		mistakes:      make(map[int]bool),
		mistakeCount:  0,

		inputs: make(map[string]textinput.Model, 1), // TODO map instead slice
	}

	t := textinput.New()
	t.CursorStyle = cursorStyle
	t.CharLimit = 32
	t.Prompt = "Words count"
	t.Placeholder = "Nickname"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	m.inputs["words"] = t

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
