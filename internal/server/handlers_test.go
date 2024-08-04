package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

type TestStorage struct{}

func (ts *TestStorage) Get(mType common.MetricType, name string) (common.Metrics, error) {
	delta := int64(111)
	value := float64(222.444)

	metric := common.Metrics{
		ID: name,
	}

	switch mType {
	case common.CounterMetric:
		metric.MType = string(common.CounterMetric)
		metric.Delta = &delta
	case common.GaugeMetric:
		metric.MType = string(common.GaugeMetric)
		metric.Value = &value
	}

	if name == "error_metric" {
		return common.Metrics{}, errors.New("Get() error")
	}

	return metric, nil
}
func (ts *TestStorage) Set(common.Metrics) error {
	return errors.New("Set() error")
}

func (ts *TestStorage) GetData() (map[string]common.Metrics, error) {
	return nil, errors.New("GetData() error")
}
func (ts *TestStorage) SetData(map[string]common.Metrics) error {
	return errors.New("SetData() error")
}

func TestHandlers_text(t *testing.T) {
	server := Server{repository: storage.NewMemRepository()}
	server.InitLogging()

	router, err := server.setupRouter()

	assert.NoError(t, err)

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "save gauge",
			request: "/update/gauge/gaugeMetric/100.35",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:    "save counter",
			request: "/update/counter/counterMetric/400",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:    "wrong request",
			request: "/test",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "save wrong type",
			request: "/update/invalid/someMetric/100",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "invalid gauge value",
			request: "/update/gauge/gaugeMetric/invalid",
			method:  http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "invalid counter value",
			request: "/update/counter/counterMetric/invalid",
			method:  http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "get valid gauge",
			request: "/value/gauge/gaugeMetric",
			method:  http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
				body:        "100.35",
			},
		},
		{
			name:    "get valid counter",
			request: "/value/counter/counterMetric",
			method:  http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
				body:        "400",
			},
		},
		{
			name:    "get unknown metric type",
			request: "/value/invalid/someMetric",
			method:  http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "get invalid gauge",
			request: "/value/gauge/notExists",
			method:  http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "get invalid counter",
			request: "/value/counter/notExists",
			method:  http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)

			if len(tt.want.body) > 0 {
				resultBody, err := io.ReadAll(result.Body)
				assert.Nil(t, err)

				assert.Equal(t, tt.want.body, string(resultBody))
			}
		})
	}
}

func TestHandler_errorStorage(t *testing.T) {
	server := Server{repository: &TestStorage{}}
	server.InitLogging()

	router, err := server.setupRouter()

	assert.NoError(t, err)

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "test gauge",
			request: "/update/gauge/someMetric/200",
			method:  http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "test counter",
			request: "/update/counter/someMetric/300",
			method:  http.MethodPost,
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
		})
	}
}

func BenchmarkGetMetricHandle(b *testing.B) {
	server := Server{repository: &TestStorage{}}
	server.InitLogging()

	router, err := server.setupRouter()

	assert.NoError(b, err)

	b.ResetTimer()

	b.Run("counter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodGet, "/value/counter/counterMetric", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)
		}
	})

	b.Run("gauge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodGet, "/value/gauge/gaugeMetric", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)
		}
	})
}
