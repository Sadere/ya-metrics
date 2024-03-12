package agent

import (
	"time"

	"github.com/Sadere/ya-metrics/internal/agent/config"
)

type Agent struct {
	config config.Config

	mGauge map[string]float64

	pollCount int
}

func Run() {
	agent := Agent{
		config:    config.NewConfig(),
		mGauge:    make(map[string]float64),
		pollCount: 0,
	}

	// Основной цикл работы
	for {
		// Задержка перед сбором метрик
		time.Sleep(time.Duration(agent.config.PollInterval) * time.Second)

		// Считываем метрики
		agent.Poll()

		// Задержка перед
		time.Sleep(time.Duration(agent.config.ReportInterval) * time.Second)

		agent.Report()
	}
}
