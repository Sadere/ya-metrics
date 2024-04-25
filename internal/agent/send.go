package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Sadere/ya-metrics/internal/agent/middleware"
	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-resty/resty/v2"
)

const (
	InitialInterval = time.Second // Начальный интервал для backoff
	MaxRetries      = 3           // Максимальное кол-во попыток для отправки данных
)

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

	// Используем middleware для сжатия gzip и хеширования запроса
	client.SetTransport(
		&middleware.HashRoundTripper{
			Key: []byte(a.config.CryptoKey),
			Next: &middleware.GzipRoundTripper{
				Next: http.DefaultTransport,
			},
		},
	)

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
