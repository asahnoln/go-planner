package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
	"github.com/charmbracelet/bubbles/list"
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

	List    list.Model
	curView View
	project *plan.Project
	err     error
	editing *int
}

type Item struct {
	description, duration, timeRange string
}

func (i Item) Title() string {
	return i.description
}
func (i Item) Description() string {
	return i.timeRange
}
func (i Item) FilterValue() string {
	return i.description
}

func New(p *plan.Project) Model {
	if p == nil {
		p = plan.New()
	}
	m := Model{
		project: p,
	}

	m.List = list.New(m.planItems(), list.NewDefaultDelegate(), 0, 24)
	m.List.Title = "Planner"

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
		case "ctrl+c":
			return m, tea.Quit
		}
	case ViewMsg:
		return m.switchView(View(msg))
	}

	return m.updateView(msg)
}

func (m Model) View() string {
	var b strings.Builder

	switch m.curView {
	case MainView:
		b.WriteString(m.List.View())

	case AddView:
		for _, i := range m.Inputs {
			b.WriteString(i.View() + "\n")
		}
		if m.err != nil {
			b.WriteString("Event Duration must be number in minutes!")
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

func (m Model) planItems(es ...*plan.Event) []list.Item {
	items := make([]list.Item, 0)
	for _, e := range m.project.Add(es...) {
		items = append(items, Item{
			description: e.Description,
			duration:    strconv.Itoa(int(e.Duration() / time.Minute)),
			timeRange:   e.TimeRange(),
		})
	}
	return items
}

func (m Model) updateView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.curView {
	case AddView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				return m.switchView(MainView)
			case "enter":
				if m.Inputs[0].Value() != "" && m.Inputs[1].Value() != "" {
					i, err := strconv.Atoi(m.Inputs[1].Value())
					if err != nil {
						m.err = err
						return m, nil
					}

					if m.editing != nil {
						es := m.project.Add()
						es[*m.editing].Description = m.Inputs[0].Value()
						es[*m.editing].SetDuration(time.Duration(i) * time.Minute)
						m.List.SetItems(m.planItems())
					} else {
						m.List.SetItems(m.planItems(
							plan.NewEvent(m.Inputs[0].Value(), time.Duration(i)*time.Minute),
						))
					}

					m.resetInputs()
					return m.switchView(MainView)
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

	case MainView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit
			case "a":
				return m.switchView(AddView)
			case "enter":
				if len(m.List.Items()) == 0 {
					return m, nil
				}
				item := m.List.SelectedItem().(Item)
				m.Inputs[0].SetValue(item.description)
				m.Inputs[1].SetValue(item.duration)
				i := m.List.Index()
				m.editing = &i
				return m.switchView(AddView)
			}

		}

		m.List, _ = m.List.Update(msg)
	}

	return m, nil
}

func (m Model) switchInputFocus() {
	if m.Inputs[0].Focused() {
		m.Inputs[0].Blur()
		m.Inputs[1].Focus()
	} else {
		m.Inputs[1].Blur()
		m.Inputs[0].Focus()
	}
}
