package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
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
	agent := Agent{
		Host: "localhost",
		Port: 8080,

		mu:     sync.RWMutex{},
		mGauge: make(map[string]float64),

		pollCount:      0,
		pollInterval:   2,
		reportInterval: 10,
	}

	done := make(chan bool, 1)

	// Считываем метрики
	go func() {
		for {
			agent.mu.Lock()

			agent.mGauge = agent.pollMetrics()
			agent.pollCount += 1

			agent.mu.Unlock()

			time.Sleep(time.Duration(agent.pollInterval) * time.Second)
		}
	}()

	// Сохраняем на сервере метрики из рантайма
	go func() {
		for {
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

			time.Sleep(time.Duration(agent.reportInterval) * time.Second)
		}
	}()

	<-done
}

func (a *Agent) sendMetric(metricType string, metricName string, metricValue string) error {
	baseURL := "http://" + a.Host + ":" + strconv.Itoa(a.Port)

	path := fmt.Sprintf("/update/%s/%s/%s", metricType, metricName, metricValue)

	req, err := http.NewRequest(http.MethodPost, baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("[%s] couldn't create http request", metricName)
	}

	req.Header.Set("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("[%s] couldn't make http request", metricName)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("[%s] failed to save metric, code = %d", metricName, res.StatusCode)
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
