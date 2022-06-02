package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	AddView            bool
	EventNameInput     textinput.Model
	EventDurationInput textinput.Model
	Project            *plan.Project
	Err                error
}

func NewModel() Model {
	tn := textinput.New()
	tn.Placeholder = "Event Description"
	td := textinput.New()
	td.Placeholder = "Event Duration"

	return Model{
		EventNameInput:     tn,
		EventDurationInput: td,
		Project:            plan.NewProject(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.AddView = true
		case "enter":
			if m.AddView {
				// TODO: Test Bluring (it should not capture incoming letters)
				// m.EventNameInput.Blur()
				i, err := strconv.Atoi(m.EventDurationInput.Value())
				if err != nil {
					m.Err = err
					return m, nil
				}
				m.Project.AddEvent(m.EventNameInput.Value(), time.Minute*time.Duration(i))
				m.AddView = false
			}
		// TODO: Test it out
		case "tab":
			if m.AddView {
				if m.EventNameInput.Focused() {
					m.EventDurationInput.Focus()
					m.EventNameInput.Blur()
				} else {
					m.EventDurationInput.Blur()
					m.EventNameInput.Focus()
				}
			}
		}

	}

	if m.AddView {
		m.EventNameInput, _ = m.EventNameInput.Update(msg)
		m.EventDurationInput, _ = m.EventDurationInput.Update(msg)
		if !m.EventDurationInput.Focused() && !m.EventNameInput.Focused() {
			m.EventNameInput.Focus()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	if m.AddView {
		b.WriteString(m.EventNameInput.View() + "\n")
		b.WriteString(m.EventDurationInput.View() + "\n")
		if m.Err != nil {
			b.WriteString(m.Err.Error())
		}
	} else {
		b.WriteString("Planner\n")
		b.WriteString(m.EventNameInput.Value() + " ")
		b.WriteString(m.EventDurationInput.Value() + " ")
		b.WriteString(m.Project.Events(0).TimeRange())
	}

	return b.String()
}
