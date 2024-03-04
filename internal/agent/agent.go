package agent

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/Sadere/ya-metrics/internal/server"
)

type Agent struct {
	Host string
	Port int

	mu     sync.RWMutex
	mGauge map[string]float64

	pollCount,
	pollInterval,
	reportInterval int
}

func Run() {
	addr := new(server.NetAddress)
	addr.Host = "localhost"
	addr.Port = 8080

	envAddr, hasEnvAddr := os.LookupEnv("ADDRESS")

	if hasEnvAddr {
		addr.Set(envAddr)
	} else {
		flag.Var(addr, "a", "Адрес сервера")
	}

	// Конфигурируем
	flagReportInterval := flag.Int("r", 10, "Частота опроса сервера в секундах")
	flagPollInterval := flag.Int("p", 2, "Частота сбора метрик")
	flag.Parse()

	var optPollInterval, optReportInterval int

	// Частота опроса сервера
	envReportInterval, hasEnvReportInterval := os.LookupEnv("REPORT_INTERVAL")
	if hasEnvReportInterval {
		envInt, err := strconv.Atoi(envReportInterval)
		if err != nil {
			optReportInterval = envInt
		}
	} else {
		optReportInterval = *flagReportInterval
	}

	// Частота сбора метрик
	envPollInterval, hasEnvPollInterval := os.LookupEnv("POLL_INTERVAL")
	if hasEnvPollInterval {
		envInt, err := strconv.Atoi(envPollInterval)
		if err != nil {
			optPollInterval = envInt
		}
	} else {
		optPollInterval = *flagPollInterval
	}

	agent := Agent{
		Host: addr.Host,
		Port: addr.Port,

		mu:     sync.RWMutex{},
		mGauge: make(map[string]float64),

		pollCount:      0,
		pollInterval:   optPollInterval,
		reportInterval: optReportInterval,
	}

	done := make(chan bool, 1)

	// Считываем метрики
	go func() {
		for {
			time.Sleep(time.Duration(agent.pollInterval) * time.Second)

			agent.mu.Lock()

			agent.mGauge = agent.pollMetrics()
			agent.pollCount += 1

			agent.mu.Unlock()
		}
	}()

	// Сохраняем на сервере метрики из рантайма
	go func() {
		for {
			time.Sleep(time.Duration(agent.reportInterval) * time.Second)

			agent.mu.RLock()

			for metricName, metricRaw := range agent.mGauge {
				metricValue := fmt.Sprintf("%f", metricRaw)

				err := agent.sendMetric("gauge", metricName, metricValue)
				if err != nil {
					log.Println(err.Error())
					done <- true
				}
			}

			// Сохраняем кол-во считываний
			sendPollCount := strconv.Itoa(agent.pollCount)

			if err := agent.sendMetric("counter", "PollCount", sendPollCount); err != nil {
				log.Println(err.Error())
				done <- true
			}

			// Сохраняем случайное значение
			randomValue := fmt.Sprintf("%d", rand.Intn(10000))

			if err := agent.sendMetric("gauge", "RandomValue", randomValue); err != nil {
				log.Println(err.Error())
				done <- true
			}

			agent.mu.RUnlock()
		}
	}()

	<-done
}

func (a *Agent) sendMetric(metricType string, metricName string, metricValue string) error {
	baseURL := "http://" + a.Host + ":" + strconv.Itoa(a.Port)
	client := resty.New()

	path := fmt.Sprintf("/update/%s/%s/%s", metricType, metricName, metricValue)

	result, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(baseURL + path)

	if err != nil {
		return fmt.Errorf("[%s] couldn't make http request: %s", metricName, err.Error())

	}

	if result.StatusCode() != http.StatusOK {
		return fmt.Errorf("[%s] failed to save metric, code = %d", metricName, result.StatusCode())
	}

	return nil
}

func (a *Agent) pollMetrics() map[string]float64 {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)
	result := make(map[string]float64)

	result["Alloc"] = float64(rtm.Alloc)
	result["BuckHashSys"] = float64(rtm.BuckHashSys)
	result["Frees"] = float64(rtm.Frees)
	result["GCCPUFraction"] = float64(rtm.GCCPUFraction)
	result["GCSys"] = float64(rtm.GCSys)
	result["HeapAlloc"] = float64(rtm.HeapAlloc)
	result["HeapIdle"] = float64(rtm.HeapIdle)
	result["HeapInuse"] = float64(rtm.HeapInuse)
	result["HeapObjects"] = float64(rtm.HeapObjects)
	result["HeapReleased"] = float64(rtm.HeapReleased)
	result["HeapSys"] = float64(rtm.HeapSys)
	result["LastGC"] = float64(rtm.LastGC)
	result["Lookups"] = float64(rtm.Lookups)
	result["MCacheInuse"] = float64(rtm.MCacheInuse)
	result["MCacheSys"] = float64(rtm.MCacheSys)
	result["MSpanInuse"] = float64(rtm.MSpanInuse)
	result["MSpanSys"] = float64(rtm.MSpanSys)
	result["Mallocs"] = float64(rtm.Mallocs)
	result["NextGC"] = float64(rtm.NextGC)
	result["NumForcedGC"] = float64(rtm.NumForcedGC)
	result["NumGC"] = float64(rtm.NumGC)
	result["OtherSys"] = float64(rtm.OtherSys)
	result["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	result["StackInuse"] = float64(rtm.StackInuse)
	result["StackSys"] = float64(rtm.StackSys)
	result["Sys"] = float64(rtm.Sys)
	result["TotalAlloc"] = float64(rtm.TotalAlloc)

	return result
}
