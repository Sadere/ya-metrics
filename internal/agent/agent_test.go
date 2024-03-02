package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_pollMetrics(t *testing.T) {
	tests := []struct {
		name string
		a    *Agent
	}{
		{
			name: "polling test",
			a:    &Agent{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.pollMetrics()

			assert.NotEmpty(t, result)
		})
	}
}
