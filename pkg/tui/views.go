package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainViewModel struct {
}

type AddViewModel struct {
	Inputs  []textinput.Model
	err     error
	editing *int
}

func (m MainViewModel) Init() tea.Cmd {
	return nil
}
func (m MainViewModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m MainViewModel) View() string {
	return ""
}

func (m AddViewModel) Init() tea.Cmd {
	return nil
}
func (m AddViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case EditingMsg:
		m.Inputs[0].SetValue(msg.item.description)
		m.Inputs[1].SetValue(msg.item.duration)
		m.editing = &msg.i
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, switchView(MainView)
		case "enter":
			if m.Inputs[0].Value() != "" && m.Inputs[1].Value() != "" {
				i, err := strconv.Atoi(m.Inputs[1].Value())
				if err != nil {
					m.err = err
					return m, nil
				}

				if m.editing != nil {
					return m, updatePlan(*m.editing, m.Inputs[0].Value(), time.Duration(i)*time.Minute)
				} else {
					return m, insertPlan(plan.NewEvent(m.Inputs[0].Value(), time.Duration(i)*time.Minute))
				}

				m.resetInputs()
				return m, switchView(MainView)
			}

			// TODO: Actually only one way tested
			m.switchInputFocus()

		case "tab", "shift+tab":
			m.switchInputFocus()
		}
	}

	for i := range inputs {
		m.Inputs[i], _ = m.Inputs[i].Update(msg)
	}
	return m, nil
}
func (m AddViewModel) View() string {
	var b strings.Builder
	for _, i := range m.Inputs {
		b.WriteString(i.View() + "\n")
	}
	if m.err != nil {
		b.WriteString("Event Duration must be number in minutes!")
	}
	return b.String()
}

func (m AddViewModel) switchInputFocus() {
	if m.Inputs[0].Focused() {
		m.Inputs[0].Blur()
		m.Inputs[1].Focus()
	} else {
		m.Inputs[1].Blur()
		m.Inputs[0].Focus()
	}
}

func (m AddViewModel) resetInputs() {
	m.Inputs[0].Focus()
	m.Inputs[1].Blur()
	m.Inputs[0].SetValue("")
	m.Inputs[1].SetValue("")
}
