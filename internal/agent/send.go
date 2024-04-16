package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/go-resty/resty/v2"
)

const (
	MaxRetries = 3 // Максимальное кол-во попыток для отправки данных
)

// Обертка к функции отправки данных, позволяющая контроллировать сколько попыток будет для успешной отправки
func (a *MetricAgent) trySendMetrics(metrics []common.Metrics) error {
	var err error

	timeOut := 1

	for tryCount := 0; tryCount < MaxRetries; tryCount++ {
		err = a.sendMetrics(metrics)

		if err == nil {
			break
		}

		if !errors.Is(err, ErrAgentSendFailed) {
			return err
		}

		time.Sleep(time.Duration(timeOut) * time.Second)
		timeOut += 2
	}

	return err
}

// Функция самой отправки данных метрик на сервер
func (a *MetricAgent) sendMetrics(metrics []common.Metrics) error {
	baseURL := fmt.Sprintf(
		"http://%s:%d",
		a.config.ServerAddress.Host,
		a.config.ServerAddress.Port,
	)

	client := resty.New()

	path := "/updates/"

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("couldn't create json body: %s", err.Error())
	}

	// Сжимаем тело запроса
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	_, err = gz.Write(body)
	if err != nil {
		return fmt.Errorf("couldn't write gzip data: %s", err.Error())
	}

	err = gz.Close()
	if err != nil {
		return fmt.Errorf("couldn't close gzip writer: %s", err.Error())
	}

	result, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(buf.Bytes()).
		Post(baseURL + path)

	if err != nil {
		return fmt.Errorf("%w", ErrAgentSendFailed)
	}

	if result.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to save metric, code = %d", result.StatusCode())
	}

	return nil
}
