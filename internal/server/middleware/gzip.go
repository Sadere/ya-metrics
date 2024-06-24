package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Обертка над gin.ResponseWriter позволяющая сжать тело ответа в формате gzip
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// Пишет в тело gzip данные если соблюдены условия
func (g *gzipWriter) Write(data []byte) (int, error) {
	compressableContent := []string{
		"text/html",
		"application/json",
	}

	suitableContent := false
	contentType := g.ResponseWriter.Header().Get("Content-Type")

	for _, v := range compressableContent {
		if strings.Contains(contentType, v) {
			suitableContent = true
			break
		}
	}

	// Если контент не удовлетворяет условиям для сжатия, пишем тело без сжатия
	if !suitableContent {
		return g.ResponseWriter.Write(data)
	}

	g.ResponseWriter.Header().Set("Content-Encoding", "gzip")

	return g.writer.Write(data)
}

// middleware проводит сжатие тела ответа если клиент поддерживает сжатие
func GzipCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptGzip := strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip")

		// Проверяем можно ли проводить сжатие
		if !acceptGzip {
			c.Next()
			return
		}

		// Выполняем сжатие
		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
		if err != nil {
			c.Next()
			return
		}

		defer func() {
			gz.Close()
		}()

		// Переопределяем стандартный writer от gin нашим, который упакует данные в Gzip
		c.Writer = &gzipWriter{c.Writer, gz}

		c.Next()
	}
}

// middleware проводит распаковку запроса клиента, если это указано в заголовке
func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		mustDecompress := c.Request.Header.Get("Content-Encoding") == "gzip"

		// Проверяем нужно ли проводить распаковку
		if !mustDecompress {
			c.Next()
			return
		}

		// Выполняем распаковку
		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}

		c.Request.Body = gz
		c.Next()
	}
}
