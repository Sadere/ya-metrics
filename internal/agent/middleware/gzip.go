package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// middleware позволяющий сжать тело запрсоа по алгоритму gzip
type GzipRoundTripper struct {
	Next http.RoundTripper
}

func (t *GzipRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("[gzip] couldn't read request body: %s", err.Error())
	}

	// Сжимаем тело запроса
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	_, err = gz.Write(body)
	if err != nil {
		return nil, fmt.Errorf("[gzip] couldn't write gzip data: %s", err.Error())
	}

	err = gz.Close()
	if err != nil {
		return nil, fmt.Errorf("[gzip] couldn't close gzip writer: %s", err.Error())
	}

	// Заменяем тело сжатыми данными
	r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

	// Ставим нужные заголовки
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")
	r.ContentLength = int64(len(buf.Bytes()))

	return t.Next.RoundTrip(r)
}
