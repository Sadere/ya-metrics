package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestServer_updateHandlers(t *testing.T) {
	server := Server{storage: storage.NewMemStorage()}

	router := server.setupRouter()

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
			name:    "save gauge",
			request: "/update/gauge/someMetric/100",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:    "save counter",
			request: "/update/counter/someMetric/100",
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
