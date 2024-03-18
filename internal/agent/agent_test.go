package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_pollMetrics(t *testing.T) {
	tests := []struct {
		name string
		a    *MetricAgent
	}{
		{
			name: "polling test",
			a:    &MetricAgent{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Poll()

			assert.NotEmpty(t, result)
		})
	}
}
