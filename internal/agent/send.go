package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-resty/resty/v2"
)

const (
	InitialInterval = time.Second // Начальный интервал для backoff
	MaxRetries      = 3           // Максимальное кол-во попыток для отправки данных
)

type gzipRoundTripper struct{}

func (t *gzipRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
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

	gzipReq, err := http.NewRequest(r.Method, r.URL.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("[gzip] couldn't create gzip request: %s", err.Error())
	}

	gzipReq.Header.Set("Content-Encoding", "gzip")
	gzipReq.Header.Set("Accept-Encoding", "gzip")

	return http.DefaultTransport.RoundTrip(gzipReq)
}

// Обертка к функции отправки данных, позволяющая контроллировать сколько попыток будет для успешной отправки
func (a *MetricAgent) trySendMetrics(metrics []common.Metrics) error {
	b := backoff.WithMaxRetries(
		backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(InitialInterval),
		),
		MaxRetries,
	)

	operation := func() error {
		err := a.sendMetrics(metrics)

		// Если получаем ошибку ErrAgentSendFailed то не можем продолжать попытки
		if errors.Is(err, ErrAgentSendFailed) {
			return backoff.Permanent(err)
		}

		return err
	}

	return backoff.Retry(operation, b)
}

// Функция самой отправки данных метрик на сервер
func (a *MetricAgent) sendMetrics(metrics []common.Metrics) error {
	baseURL := fmt.Sprintf(
		"http://%s:%d",
		a.config.ServerAddress.Host,
		a.config.ServerAddress.Port,
	)

	client := resty.New()

	// Используем middleware для сжатия gzip
	client.SetTransport(&gzipRoundTripper{})

	path := "/updates/"

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("couldn't create json body: %s", err.Error())
	}

	result, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(baseURL + path)

	if err != nil {
		return err
	}

	if result.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to save metric, code = %d", result.StatusCode())
	}

	return nil
}
