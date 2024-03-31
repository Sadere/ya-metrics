package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func GzipCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		compressableContent := []string{
			"text/html",
			"application/json",
		}

		acceptGzip := strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip")
		suitableContent := slices.Contains(compressableContent, c.Request.Header.Get("Content-Type"))

		// Проверяем можно ли проводить сжатие
		if !acceptGzip && !suitableContent {
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
			c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}()

		// Переопределяем стандартный writer от gin нашим, который упакует данные в Gzip
		c.Writer = &gzipWriter{c.Writer, gz}

		c.Header("Content-Encoding", "gzip")
		c.Next()
	}
}

func GzipDeCompress() gin.HandlerFunc {
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
			return
		}

		c.Request.Body = gz
		c.Next()
	}
}
