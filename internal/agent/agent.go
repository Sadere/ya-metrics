package agent

import (
	"time"

	"github.com/Sadere/ya-metrics/internal/agent/config"
)

type MetricAgent struct {
	config    config.Config
	pollCount int
}

func Run() {
	agent := MetricAgent{
		config:    config.NewConfig(),
		pollCount: 0,
	}

	// Основной цикл работы
	for {
		// Задержка перед сбором метрик
		time.Sleep(time.Duration(agent.config.PollInterval) * time.Second)

		// Считываем метрики
		gaugeMetrics := agent.Poll()

		// Задержка перед отправкой метрик на сервер
		time.Sleep(time.Duration(agent.config.ReportInterval) * time.Second)

		agent.Report(gaugeMetrics)
	}
}
