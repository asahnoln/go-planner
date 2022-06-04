package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
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
	Inputs []textinput.Model

	curView View
	project *plan.Project
}

func New(p *plan.Project) Model {
	if p == nil {
		p = plan.New()
	}
	m := Model{
		project: p,
	}

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
				i, err := strconv.Atoi(m.Inputs[1].Value())
				if err != nil {
					// TODO: Test what has to happen on err?
				}
				m.project.Add(plan.NewEvent(m.Inputs[0].Value(), time.Duration(i)*time.Minute))
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

	for _, e := range m.project.Add() {
		b.WriteString(e.Description)
		b.WriteString(" ")
		b.WriteString(e.TimeRange())
		b.WriteString("\n")
	}

	if m.curView == AddView {
		for _, i := range m.Inputs {
			b.WriteString(i.View() + "\n")
		}
		return b.String()
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
