package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/Sadere/ya-metrics/internal/common"
)

// middleware для вычисления хеша из тела запроса и передачи хеша в заголовке
type HashRoundTripper struct {
	Key  []byte
	Next http.RoundTripper
}

func (ht *HashRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if len(ht.Key) == 0 {
		return ht.Next.RoundTrip(r)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Создаем новый r.Body т.к. мы прочитали все данные
	r.Body = io.NopCloser(bytes.NewReader(body))

	h := hmac.New(sha256.New, ht.Key)
	h.Write(body)
	bodyHash := h.Sum(nil)

	r.Header.Set(common.HashHeader, hex.EncodeToString(bodyHash))

	return ht.Next.RoundTrip(r)
}
