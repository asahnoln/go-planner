package planner_test

import (
	"testing"
	"time"

	"github.com/asahnoln/go-planner"
)

func TestAddEventGetTimeRange(t *testing.T) {
	event := planner.AddEvent(5 * time.Minute)
	event2 := planner.AddEvent(10 * time.Minute)
	event3 := planner.AddEvent(15 * time.Minute)

	message := "want time range %q, got %q"
	assertSameString(t, "12:15-12:30", event3.TimeRange(), message)
	assertSameString(t, "12:00-12:05", event.TimeRange(), message)
	assertSameString(t, "12:05-12:15", event2.TimeRange(), message)
}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
