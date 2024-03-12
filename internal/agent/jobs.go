package agent

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"math/rand"

	"github.com/go-resty/resty/v2"
)

// "Процесс" сборки метрик
func (a *Agent) Poll() {
	a.mGauge = a.pollMetrics()
	a.pollCount += 1
}

// "Процесс" отправки метрик на сервер
func (a *Agent) Report() {
	for metricName, metricRaw := range a.mGauge {
		metricValue := fmt.Sprintf("%f", metricRaw)

		err := a.sendMetric("gauge", metricName, metricValue)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	// Сохраняем кол-во считываний
	sendPollCount := strconv.Itoa(a.pollCount)

	if err := a.sendMetric("counter", "PollCount", sendPollCount); err != nil {
		log.Println(err.Error())
		return
	}

	// Сохраняем случайное значение
	randomValue := fmt.Sprintf("%d", rand.Intn(10000))

	if err := a.sendMetric("gauge", "RandomValue", randomValue); err != nil {
		log.Println(err.Error())
		return
	}
}

// Функция отправки данных метрик на сервер
func (a *Agent) sendMetric(metricType string, metricName string, metricValue string) error {
	baseURL := "http://" + a.config.Host + ":" + strconv.Itoa(a.config.Port)
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

// Считывание необходимых метрик
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
