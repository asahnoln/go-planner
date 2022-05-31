package planner_test

import (
	"testing"
	"time"

	"github.com/asahnoln/go-planner"
	"github.com/stretchr/testify/assert"
)

func TestEventFullDescription(t *testing.T) {
	project := planner.NewProject()
	event := project.AddEvent("Warmup", 10*time.Minute)

	assert.Equal(t, "Warmup: 12:00-12:10", event.FullDescription(), "want full description %q, got %q")
}
