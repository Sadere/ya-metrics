package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIPValidation(t *testing.T) {
	tests := []struct {
		name     string
		subnet   string
		IP       string
		wantCode int
	}{
		{
			name:     "valid IP",
			subnet:   "192.168.1.0/24",
			IP:       "192.168.1.1",
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid IP",
			subnet:   "192.168.1.0/24",
			IP:       "100.100.1.1",
			wantCode: http.StatusForbidden,
		},
		{
			name:     "all IP allowed",
			subnet:   "0.0.0.0/0",
			IP:       "100.100.1.1",
			wantCode: http.StatusOK,
		},
		{
			name:     "malformed IP",
			subnet:   "192.168.1.0/24",
			IP:       "999.999.1.1",
			wantCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw, err := ValidateIP(tt.subnet)

			assert.NoError(t, err)

			r := gin.New()
			r.Use(mw)
			r.GET("/example", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			request := httptest.NewRequest(http.MethodGet, "/example", nil)
			request.Header.Add(common.IPHeader, tt.IP)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.wantCode, result.StatusCode)
		})
	}
}

func TestWrongSubnet(t *testing.T) {
	_, err := ValidateIP("333.444.555.666/99")

	assert.Error(t, err)
}
