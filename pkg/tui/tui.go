package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var inputs = []struct {
	placeholder string
}{
	{placeholder: "Event Description"},
	{placeholder: "Event Duration"},
}

type View int

type ViewMsg View

const (
	MainView View = iota
	AddView
)

type Model struct {
	Inputs  []textinput.Model
	curView View
	added   bool
}

func New() Model {
	m := Model{}
	for _, i := range inputs {
		t := textinput.New()
		t.Placeholder = i.placeholder
		m.Inputs = append(m.Inputs, t)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "a":
			return m.switchView(AddView)
		case "enter":
			if m.curView == AddView {
				m.added = true
				return m.switchView(MainView)
			}
		}
	case ViewMsg:
		return m.switchView(View(msg))
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Planner\n")

	if m.curView == AddView {
		for _, i := range m.Inputs {
			b.WriteString(i.View() + "\n")
		}
		return b.String()
	}

	if m.added {
		b.WriteString("Intro 5 12:00-12:05")
	}
	return b.String()
}

func (m Model) switchView(v View) (tea.Model, tea.Cmd) {
	m.curView = v
	if v == AddView {
		m.Inputs[0].Focus()
	}
	return m, nil
}
