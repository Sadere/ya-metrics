package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
)

type bodyHashWriter struct {
	gin.ResponseWriter

	statusCode int
	h          http.Header
	body       *bytes.Buffer
}

func (w *bodyHashWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *bodyHashWriter) Header() http.Header {
	return w.h
}

func (w *bodyHashWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *bodyHashWriter) Done() {
	originalHeader := w.ResponseWriter.Header()
	for key, value := range w.Header() {
		originalHeader[key] = value
	}

	w.ResponseWriter.WriteHeader(w.statusCode)

	_, err := w.ResponseWriter.Write(w.body.Bytes())
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Сверяем хеш запроса
func ValidateHash(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		providedHash := c.Request.Header.Get(common.HashHeader)

		if providedHash == "" {
			c.Next()
			return
		}

		// Читаем тело запроса и считаем хеш из ключа
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "unexpected error: %s", err.Error())
			c.Abort()
			return
		}

		// Возвращаем на место тело запроса
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		h := hmac.New(sha256.New, []byte(key))
		h.Write(body)

		ourHash := h.Sum(nil)

		theirHash, err := hex.DecodeString(providedHash)
		if err != nil {
			c.String(http.StatusBadRequest, "failed to decode hash")
			c.Abort()
			return
		}

		if !hmac.Equal(ourHash, theirHash) {
			c.String(http.StatusBadRequest, "hash mismatch")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Считаем хеш ответа от сервера
func HashResponse(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(key) == 0 {
			return
		}

		// Подменяем writer
		hashWriter := &bodyHashWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
			h:              make(http.Header),
		}
		c.Writer = hashWriter

		defer hashWriter.Done()

		// Делаем что-то
		c.Next()

		// Хешируем ответ
		h := hmac.New(sha256.New, []byte(key))
		h.Write(hashWriter.body.Bytes())

		responseHash := h.Sum(nil)

		// Ставим заголовок
		c.Header(common.HashHeader, hex.EncodeToString(responseHash))
	}
}
