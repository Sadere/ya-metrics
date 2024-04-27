package agent

import (
	"testing"
	"time"

	"github.com/Sadere/ya-metrics/internal/agent/config"
	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestAgent_pollMetrics(t *testing.T) {
	tests := []struct {
		name string
		a    *MetricAgent
	}{
		{
			name: "polling test",
			a: &MetricAgent{
				config: config.Config{
					ReportInterval: 1,
					PollInterval:   1,
				},
				workChan: make(chan []common.Metrics),
				doneChan: make(chan struct{}, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				for {
					select {
					case <-tt.a.doneChan:
						return
					case metrics := <-tt.a.workChan:
						assert.NotEmpty(t, metrics)
					}
				}
			}()

			tt.a.PollRuntime()

			time.Sleep(time.Second * 2)

			tt.a.doneChan <- struct{}{}
		})
	}
}
