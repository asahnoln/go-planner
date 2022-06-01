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
			m.EventNameInput.Focus()
			m.AddView = true
		case "enter":
			if m.AddView {
				m.AddView = false
			}
		}

	}

	var cmd tea.Cmd
	m.EventNameInput, cmd = m.EventNameInput.Update(msg)
	return m, cmd
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
