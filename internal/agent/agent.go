package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sadere/ya-metrics/internal/agent/config"
	"github.com/Sadere/ya-metrics/internal/common"
)

// Информация о сборке
var (
	buildVersion string
	buildDate string
	buildCommit string
)

// Главная структура агента
type MetricAgent struct {
	config    config.Config
	pollCount int64
	workChan  chan []common.Metrics
	doneChan  chan struct{}
}

func (a *MetricAgent) worker(id int) {
	for {
		select {
		case <-a.doneChan:
			log.Printf("worker #%d shutdown...\n", id+1)
			return
		case metrics := <-a.workChan:

			// Задержка перед отправкой метрик на сервер
			time.Sleep(time.Duration(a.config.ReportInterval) * time.Second)

			err := a.trySendMetrics(metrics)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

// Основной метод агента, запускает расчет метрик и отправку их на сервер
func Run() {
	// Выводим информацию о сборке
	fmt.Print(common.BuildInfo(buildVersion, buildDate, buildCommit))

	agent := MetricAgent{
		config:    config.NewConfig(),
		pollCount: 0,
		workChan:  make(chan []common.Metrics),
	}

	agent.doneChan = make(chan struct{}, agent.config.RateLimit)

	for i := 0; i < agent.config.RateLimit; i++ {
		go agent.worker(i)
	}

	go func() {
		for {
			// Задержка перед сбором метрик
			time.Sleep(time.Duration(agent.config.PollInterval) * time.Second)

			// Считываем метрики
			agent.PollRuntime()
			agent.PollPS()
		}
	}()

	// Ловим сигналы отключения агента
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Закрываем всех воркеров
	for i := 0; i < agent.config.RateLimit; i++ {
		agent.doneChan <- struct{}{}
	}

	// Ждем пока воркеры завершатся
	time.Sleep(time.Second)
}
