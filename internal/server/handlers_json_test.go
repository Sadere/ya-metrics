package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_updateJSON(t *testing.T) {
	server := Server{repository: storage.NewMemRepository()}
	server.InitLogging()

	router := server.setupRouter()

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name        string
		requestUri  string
		requestBody []byte
		want        want
	}{
		{
			name:        "invalid json input",
			requestUri:  "/update/",
			requestBody: []byte(`not json`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:        "save gauge",
			requestUri:  "/update/",
			requestBody: []byte(`{"id":"gaugeMetric","type":"gauge","value":100.66}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:        "save counter",
			requestUri:  "/update/",
			requestBody: []byte(`{"id":"counterMetric","type":"counter","delta":200}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:        "save invalid type",
			requestUri:  "/update/",
			requestBody: []byte(`{"id":"invalidMetric","type":"invalid","delta":200}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:        "get gauge",
			requestUri:  "/value/",
			requestBody: []byte(`{"id":"gaugeMetric","type":"gauge"}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				body: `{"id":"gaugeMetric","type":"gauge","value":100.66}`,
			},
		},
		{
			name:        "get counter",
			requestUri:  "/value/",
			requestBody: []byte(`{"id":"counterMetric","type":"counter"}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				body: `{"id":"counterMetric","type":"counter","delta":200}`,
			},
		},
		{
			name:        "get invalid",
			requestUri:  "/value/",
			requestBody: []byte(`{"id":"invalidMetric","type":"invalid"}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:        "value invalid json",
			requestUri:  "/value/",
			requestBody: []byte(`not json`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:        "batch save gauge",
			requestUri:  "/updates/",
			requestBody: []byte(`[{"id":"gaugeMetric","type":"gauge","value":100.66}]`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:        "batch save counter",
			requestUri:  "/updates/",
			requestBody: []byte(`[{"id":"counterMetric","type":"counter","delta":200}]`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:        "batch save invalid type",
			requestUri:  "/updates/",
			requestBody: []byte(`[{"id":"invalidMetric","type":"invalid"}]`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:        "batch invalid json",
			requestUri:  "/updates/",
			requestBody: []byte(`not json`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestUri, bytes.NewReader(tt.requestBody))
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

func TestHandlerJSON_errorStorage(t *testing.T) {
	server := Server{repository: &TestStorage{}}
	server.InitLogging()

	router := server.setupRouter()

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name        string
		requestUri  string
		requestBody []byte
		want        want
	}{
		{
			name:        "error update handler",
			requestUri:  "/update/",
			requestBody: []byte(`{"id":"gaugeMetric","type":"gauge","value":100.66}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:        "error get handler",
			requestUri:  "/value/",
			requestBody: []byte(`{"id":"gaugeMetric","type":"gauge"}`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:        "error batch update handler",
			requestUri:  "/updates/",
			requestBody: []byte(`[{"id":"gaugeMetric","type":"gauge"}]`),
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestUri, bytes.NewReader(tt.requestBody))
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
