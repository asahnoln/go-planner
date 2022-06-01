package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	AddView        bool
	EventNameInput textinput.Model
}

func NewModel() Model {
	e := textinput.New()
	e.Placeholder = "Event Description"
	return Model{
		EventNameInput: e,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.AddView = true
		case "enter":
			if m.AddView {
				// TODO: Test Bluring (it should not capture incoming letters)
				// m.EventNameInput.Blur()
				m.AddView = false
			}
		}

	}

	if m.AddView {
		m.EventNameInput, _ = m.EventNameInput.Update(msg)
		m.EventNameInput.Focus()
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	if m.AddView {
		b.WriteString(m.EventNameInput.View())
	} else {
		b.WriteString("Planner ")
		b.WriteString(m.EventNameInput.Value())
	}

	return b.String()
}
