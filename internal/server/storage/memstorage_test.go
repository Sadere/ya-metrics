package storage

import (
	"testing"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestSetGetCounter(t *testing.T) {
	m := NewMemRepository()

	delta := int64(1111)
	err := m.Set(common.Metrics{
		ID:    "test_counter",
		MType: string(common.CounterMetric),
		Delta: &delta,
	})

	assert.Nil(t, err)

	// Check value
	val, err := m.Get(common.CounterMetric, "test_counter")
	assert.Nil(t, err)
	assert.Equal(t, delta, *val.Delta)

}
func TestSetGetGauge(t *testing.T) {
	m := NewMemRepository()

	value := float64(1111)
	err := m.Set(common.Metrics{
		ID:    "test_gauge",
		MType: string(common.GaugeMetric),
		Value: &value,
	})

	assert.Nil(t, err)

	// Check value
	val, err := m.Get(common.GaugeMetric, "test_gauge")
	assert.Nil(t, err)
	assert.Equal(t, value, *val.Value)
}

func TestSetUnknown(t *testing.T) {
	m := NewMemRepository()

	err := m.Set(common.Metrics{
		ID:    "test",
		MType: "unknown",
	})

	assert.NotNil(t, err)
}

func TestGetUnknown(t *testing.T) {
	m := NewMemRepository()

	_, err := m.Get("unkown", "test")

	assert.NotNil(t, err)
}

func TestSetEmpty(t *testing.T) {
	m := NewMemRepository()

	delta := int64(1111)
	err := m.Set(common.Metrics{
		ID:    "",
		MType: string(common.CounterMetric),
		Delta: &delta,
	})

	assert.NotNil(t, err)
}

func TestGetEmpty(t *testing.T) {
	m := NewMemRepository()

	_, err := m.Get(common.CounterMetric, "")

	assert.NotNil(t, err)
}

func TestGetCounterInvalid(t *testing.T) {
	m := NewMemRepository()

	_, err := m.Get(common.CounterMetric, "invalid")

	assert.NotNil(t, err)
}

func TestGetGaugeInvalid(t *testing.T) {
	m := NewMemRepository()

	_, err := m.Get(common.GaugeMetric, "invalid")

	assert.NotNil(t, err)
}

func TestSetGetData(t *testing.T) {
	m := NewMemRepository()

	value := float64(1111)
	delta := int64(2222)

	err := m.SetData(map[string]common.Metrics{
		"test_gauge": {
			ID:    "test_gauge",
			MType: string(common.GaugeMetric),
			Value: &value,
		},
		"test_counter": {
			ID:    "test_counter",
			MType: string(common.CounterMetric),
			Delta: &delta,
		},
	})

	assert.Nil(t, err)

	// Check value
	data, err := m.GetData()
	assert.Nil(t, err)

	if assert.Contains(t, data, "test_gauge") {
		assert.Equal(t, value, *data["test_gauge"].Value)
	}

	if assert.Contains(t, data, "test_counter") {
		assert.Equal(t, delta, *data["test_counter"].Delta)
	}
}
