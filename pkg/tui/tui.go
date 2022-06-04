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
	m.Inputs[0].Focus()

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
			if m.curView == MainView {
				return m.switchView(AddView)
			}
		}
	case ViewMsg:
		return m.switchView(View(msg))
	}

	return m.updateView(msg)
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

	return m, nil
}

func (m Model) resetInputs() {
	m.Inputs[0].Focus()
	m.Inputs[1].Blur()
	m.Inputs[0].SetValue("")
	m.Inputs[1].SetValue("")
}

func (m Model) updateView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.curView {
	case AddView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				i, err := strconv.Atoi(m.Inputs[1].Value())
				if err != nil {
					// TODO: Test what has to happen on err?
				}
				m.project.Add(plan.NewEvent(m.Inputs[0].Value(), time.Duration(i)*time.Minute))
				m.resetInputs()
				return m.switchView(MainView)
			case "tab":
				if m.Inputs[0].Focused() {
					m.Inputs[0].Blur()
					m.Inputs[1].Focus()
				} else {
					m.Inputs[1].Blur()
					m.Inputs[0].Focus()
				}
			}
		}

		m.Inputs[0], _ = m.Inputs[0].Update(msg)
		m.Inputs[1], _ = m.Inputs[1].Update(msg)
	}

	return m, nil
}
