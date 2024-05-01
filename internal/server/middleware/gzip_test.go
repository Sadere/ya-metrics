package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	exampleGzip := []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x4a, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x55, 0x28, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x5, 0x4, 0x0, 0x0, 0xff, 0xff, 0x9e, 0xc5, 0x99, 0xae, 0x10, 0x0, 0x0, 0x0}

	r := gin.New()

	r.Use(GzipDecompress())
	r.Use(GzipCompress())
	r.Use(gin.Recovery())

	r.GET("/example", func(c *gin.Context) {
		c.String(http.StatusOK, "example response")
	})

	type want struct {
		contentEncoding string
		statusCode      int
	}
	tests := []struct {
		name    string
		headers map[string]string
		want    want
		body    []byte
	}{
		{
			name: "accept encoding",
			headers: map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "text/html",
			},
			want: want{
				contentEncoding: "gzip",
				statusCode:      http.StatusOK,
			},
		},
		{
			name: "non gzip encoding",
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				contentEncoding: "",
				statusCode:      http.StatusOK,
			},
		},
		{
			name: "gzip content",
			headers: map[string]string{
				"Accept-Encoding":  "gzip",
				"Content-Type":     "text/html",
				"Content-Encoding": "gzip",
			},
			body: exampleGzip,
			want: want{
				contentEncoding: "gzip",
				statusCode:      http.StatusOK,
			},
		},
		{
			name: "invalid gzip content",
			headers: map[string]string{
				"Accept-Encoding":  "gzip",
				"Content-Type":     "text/html",
				"Content-Encoding": "gzip",
			},
			body: []byte{0xde, 0xad, 0xbe, 0xef},
			want: want{
				contentEncoding: "",
				statusCode:      http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(tt.body)

			request := httptest.NewRequest(http.MethodGet, "/example", buf)

			for headerName, headerValue := range tt.headers {
				request.Header.Add(headerName, headerValue)
			}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Encoding"), tt.want.contentEncoding)
		})
	}
}
