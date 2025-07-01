package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Sadere/ya-metrics/internal/agent/config"
	"github.com/Sadere/ya-metrics/internal/agent/middleware"
	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-resty/resty/v2"
)

const (
	InitialInterval = time.Second // Начальный интервал для backoff
	MaxRetries      = 3           // Максимальное кол-во попыток для отправки данных
)

var (
	ErrAgentSendFailed = errors.New("agent couldn't transfer data to server")
)

// Транспорт для отправки метрик по HTTP
type HTTPMetricTransport struct {
	config config.Config
}

func NewHTTPMetricTransport(cfg config.Config) *HTTPMetricTransport {
	return &HTTPMetricTransport{
		config: cfg,
	}
}

// Обертка к функции отправки данных, позволяющая контроллировать сколько попыток будет для успешной отправки
func (t *HTTPMetricTransport) SendMetrics(metrics []common.Metrics) error {
	b := backoff.WithMaxRetries(
		backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(InitialInterval),
		),
		MaxRetries,
	)

	operation := func() error {
		err := t.sendMetrics(metrics)

		// Если получаем ошибку ErrAgentSendFailed то не можем продолжать попытки
		if errors.Is(err, ErrAgentSendFailed) {
			return backoff.Permanent(err)
		}

		return err
	}

	return backoff.Retry(operation, b)
}

// Функция самой отправки данных метрик на сервер
func (t *HTTPMetricTransport) sendMetrics(metrics []common.Metrics) error {
	baseURL := fmt.Sprintf(
		"http://%s:%d",
		t.config.ServerAddress.Host,
		t.config.ServerAddress.Port,
	)

	client := resty.New()

	// Настраиваем middleware
	transport := &middleware.HashRoundTripper{
		Key: []byte(t.config.HashKey),
	}
	gzipTransport := &middleware.GzipRoundTripper{
		Next: http.DefaultTransport,
	}

	// middleware для шифрования
	if len(t.config.PubKeyFilePath) > 0 {
		transport.Next = &middleware.CryptoRoundTripper{
			KeyFilePath: t.config.PubKeyFilePath,
			Next:        gzipTransport,
		}
	} else {
		transport.Next = gzipTransport
	}

	client.SetTransport(transport)

	path := "/updates/"

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("couldn't create json body: %s", err.Error())
	}

	result, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader(common.IPHeader, t.config.HostAddress).
		SetBody(body).
		Post(baseURL + path)

	if err != nil {
		return ErrAgentSendFailed
	}

	if result.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to save metric, code = %d", result.StatusCode())
	}

	return nil
}