package tui_test

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
	"github.com/asahnoln/go-planner/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	m := tui.New(nil)
	assert.Contains(t, m.View(), "Planner")
}

func TestInitWithPlan(t *testing.T) {
	p := plan.New()
	_ = p.Add(plan.NewEvent("Prepare", 10*time.Minute),
		plan.NewEvent("Ignite", 5*time.Minute),
	)

	m := tui.New(p)
	assert.Contains(t, m.View(), "Planner")
	assert.Contains(t, m.View(), "Ignite")
}

func TestQuit(t *testing.T) {
	tests := []struct {
		name    string
		keyType tea.KeyType
		runes   []rune
	}{
		{"q", tea.KeyRunes, []rune{'q'}},
		{"ctrl+c", tea.KeyCtrlC, nil},
		// TODO: Esc might be not good for subviews
		{"esc", tea.KeyEsc, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tui.New(nil)
			_, c := m.Update(tea.KeyMsg(tea.Key{
				Type:  tt.keyType,
				Runes: tt.runes,
			}))

			qn := runtime.FuncForPC(reflect.ValueOf(tea.Quit).Pointer()).Name()
			cn := runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()
			assert.Equal(t, qn, cn, "want quitMsg")
		})
	}
}

// TODO: When pressed a - it's sent to the input
func TestSwitchToAddEvent(t *testing.T) {
	var m tea.Model = tui.New(nil)
	m, _ = m.Update(tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'a'},
	}))
	n := m.(tui.Model)
	v := m.View()
	require.Contains(t, v, "> ", "want text input")
	assert.Contains(t, v, "Description", "want text input placeholder")
	assert.Contains(t, v, "Duration", "want time input placeholder")
	assert.True(t, n.Inputs[0].Focused(), "want text input focused")
}

func TestAddEvent(t *testing.T) {
	var m tea.Model = tui.New(nil)
	n := m.(tui.Model)
	n.Inputs[0].SetValue("Intro")
	n.Inputs[1].SetValue("5")

	t.Run("enter should not work on main page", func(t *testing.T) {
		m, _ = n.Update(tui.ViewMsg(tui.MainView))
		m, _ = m.Update(tea.KeyMsg(tea.Key{
			Type: tea.KeyEnter,
		}))
		v := m.View()
		assert.NotContains(t, v, "Intro", "want event name on the page")
	})

	t.Run("usual add", func(t *testing.T) {
		m, _ = n.Update(tui.ViewMsg(tui.AddView))
		m, _ = m.Update(tea.KeyMsg(tea.Key{
			Type: tea.KeyEnter,
		}))
		v := m.View()
		assert.Contains(t, v, "Planner", "want return to main page")
		assert.Contains(t, v, "Intro", "want event name on the page")
		assert.Contains(t, v, "5", "want event duration on the page")
		assert.Contains(t, v, "12:00-12:05", "want event timings on the page")
	})
}

func TestSelectInputs(t *testing.T) {
	var m tea.Model = tui.New(nil)

	t.Run("tab should not work with inputs on main page", func(t *testing.T) {
		m, _ = m.Update(tui.ViewMsg(tui.MainView))
		m, _ = m.Update(tea.KeyMsg(tea.Key{
			Type: tea.KeyTab,
		}))
		n := m.(tui.Model)
		assert.True(t, n.Inputs[0].Focused(), "want first input still focused")
		assert.False(t, n.Inputs[1].Focused(), "want second input still blured")
	})

	t.Run("usual switch focus", func(t *testing.T) {
		m, _ = m.Update(tui.ViewMsg(tui.AddView))
		m, _ = m.Update(tea.KeyMsg(tea.Key{
			Type: tea.KeyTab,
		}))
		n := m.(tui.Model)
		assert.False(t, n.Inputs[0].Focused(), "want first input blured")
		assert.True(t, n.Inputs[1].Focused(), "want first input focused")

		m, _ = n.Update(tea.KeyMsg(tea.Key{
			Type: tea.KeyTab,
		}))
		n = m.(tui.Model)
		assert.True(t, n.Inputs[0].Focused(), "want first input focused back")
		assert.False(t, n.Inputs[1].Focused(), "want second input blured back")
	})
}

func TestInput(t *testing.T) {
	var m tea.Model = tui.New(nil)
	m, _ = m.Update(tui.ViewMsg(tui.AddView))
	m, _ = m.Update(tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'a'},
	}))
	n := m.(tui.Model)
	assert.Equal(t, "a", n.Inputs[0].Value(), "want first input have entered letter 'a'")

	m, _ = n.Update(tea.KeyMsg(tea.Key{
		Type: tea.KeyTab,
	}))
	m, _ = m.Update(tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'a'},
	}))
	n = m.(tui.Model)
	assert.Equal(t, "a", n.Inputs[1].Value(), "want second input have entered letters 'a'")
}

func TestResetFocusAfterAdding(t *testing.T) {
	var m tea.Model = tui.New(nil)
	m, _ = m.Update(tui.ViewMsg(tui.AddView))

	n := m.(tui.Model)
	n.Inputs[1].Focus()
	n.Inputs[0].SetValue("test value 1")
	n.Inputs[1].SetValue("test value 2")

	m, _ = n.Update(tea.KeyMsg(tea.Key{
		Type: tea.KeyEnter,
	}))

	n = m.(tui.Model)
	assert.True(t, n.Inputs[0].Focused(), "want inputs reset: first focus")
	assert.False(t, n.Inputs[1].Focused(), "want inputs reset: second blur")
	assert.Empty(t, n.Inputs[0].Value(), "want inputs reset: first empty")
	assert.Empty(t, n.Inputs[1].Value(), "want inputs reset: second empty")
}
