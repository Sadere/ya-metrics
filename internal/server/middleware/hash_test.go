package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	exampleText := "example text"
	hashKey := "example key"

	h := hmac.New(sha256.New, []byte(hashKey))
	h.Write([]byte(exampleText))
	exampleHash := hex.EncodeToString(h.Sum(nil))

	r := gin.New()

	r.Use(ValidateHash(hashKey))
	r.Use(HashResponse(hashKey))
	r.Use(gin.Recovery())

	r.GET("/example", func(c *gin.Context) {
		c.String(http.StatusOK, exampleText)
	})

	type want struct {
		hash       string
		statusCode int
	}
	tests := []struct {
		name string
		hash string
		want want
		body []byte
	}{
		{
			name: "correct hash",
			hash: exampleHash,
			body: []byte(exampleText),
			want: want{
				hash:       exampleHash,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "invalid header value",
			hash: "invalid hash header",
			body: []byte(exampleText),
			want: want{
				hash:       "",
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "wrong hash",
			hash: "aaaa",
			body: []byte(exampleText),
			want: want{
				hash:       "",
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "no hash",
			hash: "",
			body: []byte(exampleText),
			want: want{
				hash:       exampleHash,
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(tt.body)

			request := httptest.NewRequest(http.MethodGet, "/example", buf)
			request.Header.Add(common.HashHeader, tt.hash)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.hash, result.Header.Get(common.HashHeader))
		})
	}
}
