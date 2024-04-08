package middleware

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestLogger(t *testing.T) {
	var buffer bytes.Buffer

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buffer)

	memLogger := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel),
	)

	// arrange - init Gin to use the structured logger middleware
	r := gin.New()
	r.Use(Logger(memLogger))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})
	r.GET("/force500", func(c *gin.Context) { panic("forced panic") })

	// act & assert
	PerformRequest(r, "GET", "/example?a=100")
	writer.Flush()
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	PerformRequest(r, "GET", "/notfound")
	writer.Flush()
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")

	buffer.Reset()
	PerformRequest(r, "GET", "/force500")
	writer.Flush()
	assert.Contains(t, buffer.String(), "500")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/force500")
	assert.Contains(t, buffer.String(), "ERROR")
}
