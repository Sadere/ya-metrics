package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sadere/ya-metrics/internal/agent/config"
	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func inputData() []common.Metrics {
	return []common.Metrics{
		{
			ID:    "gaugeMetric",
			MType: string(common.GaugeMetric),
		},
		{
			ID:    "counterMetric",
			MType: string(common.CounterMetric),
		},
	}
}

func TestSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	addr := common.NetAddress{}
	err := addr.Set(strings.Replace(ts.URL, "http://", "", 1))
	assert.Nil(t, err)

	agent := MetricAgent{
		config: config.Config{
			ServerAddress: addr,
		},
	}

	err = agent.trySendMetrics(inputData())

	assert.Nil(t, err)
}

func TestSendError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	addr := common.NetAddress{}
	err := addr.Set(strings.Replace(ts.URL, "http://", "", 1))
	assert.Nil(t, err)

	agent := MetricAgent{
		config: config.Config{
			ServerAddress: addr,
		},
	}

	err = agent.trySendMetrics(inputData())

	assert.NotNil(t, err)
}
