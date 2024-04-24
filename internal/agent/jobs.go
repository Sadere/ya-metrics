package agent

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/Sadere/ya-metrics/internal/common"
)

// "Процесс" сборки метрик
func (a *MetricAgent) Poll() map[string]float64 {
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

	// Увеличиваем счетчик
	a.pollCount += 1

	return result
}

// "Процесс" отправки метрик на сервер
func (a *MetricAgent) Report(gaugeMetrics map[string]float64) {
	var metricsToSend []common.Metrics

	for metricName, metricValue := range gaugeMetrics {
		m := common.Metrics{
			ID:    metricName,
			MType: string(common.GaugeMetric),
			Value: &metricValue,
		}

		metricsToSend = append(metricsToSend, m)
	}

	err := a.trySendMetrics(metricsToSend)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Сохраняем кол-во считываний
	pollCount := int64(a.pollCount)
	pollCountMetric := common.Metrics{
		ID:    "PollCount",
		MType: string(common.CounterMetric),
		Delta: &pollCount,
	}

	if err := a.trySendMetrics([]common.Metrics{pollCountMetric}); err != nil {
		log.Println(err.Error())
		return
	}

	// Сохраняем случайное значение
	randomValue := float64(rand.Intn(10000))
	randomValueMetric := common.Metrics{
		ID:    "RandomValue",
		MType: string(common.GaugeMetric),
		Value: &randomValue,
	}

	if err := a.trySendMetrics([]common.Metrics{randomValueMetric}); err != nil {
		log.Println(err.Error())
		return
	}
}
