package tui_test

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/asahnoln/go-planner/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	m := tui.New()
	assert.Contains(t, m.View(), "Planner")
}

// TODO: When pressed a - it's sent to the input
// func TestSwitchToAddEvent(t *testing.T) {
// 	var m tea.Model = tui.NewModel()
// 	m, _ = m.Update(tea.KeyMsg(tea.Key{
// 		Type:  tea.KeyRunes,
// 		Runes: []rune{'a'},
// 	}))
// 	n := m.(tui.Model)
// 	v := m.View()
// 	require.Contains(t, v, "> ", "want text input")
// 	assert.Contains(t, v, "Description", "want text input placeholder")
// 	assert.Contains(t, v, "Duration", "want time input placeholder")
// 	assert.True(t, n.EventNameInput.Focused(), "want text input focused")

// }

// func TestAddEvent(t *testing.T) {
// 	var m tea.Model = tui.NewModel()
// 	n := m.(tui.Model)
// 	n.SwitchView(tui.AddView)
// 	n.EventNameInput.SetValue("Intro")
// 	n.EventDurationInput.SetValue("5")

// 	m, _ = n.Update(tea.KeyMsg(tea.Key{
// 		Type: tea.KeyEnter,
// 	}))
// 	v := m.View()
// 	assert.Contains(t, v, "Planner", "want return to main page")
// 	assert.Contains(t, v, "Intro", "want event name on the page")
// 	assert.Contains(t, v, "5", "want event duration on the page")
// 	assert.Contains(t, v, "12:00-12:05", "want event timings on the page")
// }

// func TestSelectInputs(t *testing.T) {
// 	var m tea.Model = tui.NewModel()
// 	n := m.(tui.Model)
// 	n.SwitchView(tui.AddView)

// 	m, _ = n.Update(tea.KeyMsg(tea.Key{
// 		Type: tea.KeyTab,
// 	}))
// 	assert.False(t, n.EventNameInput.Focused(), "want name input blured")
// 	// assert.True(t, n.EventDurationInput.Focused(), "want duration input focused")
// }

func TestQuit(t *testing.T) {
	tests := []struct {
		name    string
		keyType tea.KeyType
		runes   []rune
	}{
		{"q", tea.KeyRunes, []rune{'q'}},
		{"ctrl+c", tea.KeyCtrlC, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tui.New()
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
