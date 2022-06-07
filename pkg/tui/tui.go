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

type UpdatePlanMsg struct {
	i           int
	description string
	d           time.Duration
}

type InsertPlanMsg struct {
	e *plan.Event
}

type EditingMsg struct {
	i    int
	item Item
}

func switchView(v View) tea.Cmd {
	return func() tea.Msg {
		return ViewMsg(v)
	}
}

func updatePlan(i int, description string, d time.Duration) tea.Cmd {
	return func() tea.Msg {
		return UpdatePlanMsg{
			i, description, d,
		}
	}
}

func insertPlan(e *plan.Event) tea.Cmd {
	return func() tea.Msg {
		return InsertPlanMsg{e}
	}
}

func setEditing(i int, item Item) tea.Cmd {
	return func() tea.Msg {
		return EditingMsg{i, item}
	}
}

const (
	MainView View = iota
	AddView
)

type Model struct {
	Views map[View]tea.Model

	List    list.Model
	curView View
	project *plan.Project
	err     error
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
		Views: map[View]tea.Model{
			MainView: MainViewModel{},
		},
	}

	m.List = list.New(m.planItems(), list.NewDefaultDelegate(), 0, 24)
	m.List.Title = "Planner"

	am := AddViewModel{}
	for _, i := range inputs {
		t := textinput.New()
		t.Placeholder = i.placeholder
		am.Inputs = append(am.Inputs, t)
	}
	am.Inputs[0].Focus()
	m.Views[AddView] = am

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
	case UpdatePlanMsg:
		es := m.project.Add()
		es[msg.i].Description = msg.description
		es[msg.i].SetDuration(msg.d)
		m.List.SetItems(m.planItems())

		return m.switchView(MainView)
	case InsertPlanMsg:
		m.List.SetItems(m.planItems(
			msg.e,
		))

		return m.switchView(MainView)
	}

	return m.updateView(msg)
}

func (m Model) View() string {
	var b strings.Builder

	switch m.curView {
	case MainView:
		b.WriteString(m.List.View())

	case AddView:
		return m.Views[AddView].View()
	}

	return b.String()
}

func (m Model) switchView(v View) (tea.Model, tea.Cmd) {
	m.curView = v

	return m, nil
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
		return m.Views[AddView].Update(msg)
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
				i := m.List.Index()
				item := m.List.SelectedItem().(Item)

				tea.Batch(setEditing(i, item), switchView(AddView))
				return m.switchView(AddView)
			}

		}

		m.List, _ = m.List.Update(msg)
	}

	return m, nil
}
