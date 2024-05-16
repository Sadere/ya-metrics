package agent

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync/atomic"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Собираем метрики из пакета runtime
func (a *MetricAgent) PollRuntime() {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)
	rtmMetrics := make(map[string]float64)

	rtmMetrics["Alloc"] = float64(rtm.Alloc)
	rtmMetrics["BuckHashSys"] = float64(rtm.BuckHashSys)
	rtmMetrics["Frees"] = float64(rtm.Frees)
	rtmMetrics["GCCPUFraction"] = float64(rtm.GCCPUFraction)
	rtmMetrics["GCSys"] = float64(rtm.GCSys)
	rtmMetrics["HeapAlloc"] = float64(rtm.HeapAlloc)
	rtmMetrics["HeapIdle"] = float64(rtm.HeapIdle)
	rtmMetrics["HeapInuse"] = float64(rtm.HeapInuse)
	rtmMetrics["HeapObjects"] = float64(rtm.HeapObjects)
	rtmMetrics["HeapReleased"] = float64(rtm.HeapReleased)
	rtmMetrics["HeapSys"] = float64(rtm.HeapSys)
	rtmMetrics["LastGC"] = float64(rtm.LastGC)
	rtmMetrics["Lookups"] = float64(rtm.Lookups)
	rtmMetrics["MCacheInuse"] = float64(rtm.MCacheInuse)
	rtmMetrics["MCacheSys"] = float64(rtm.MCacheSys)
	rtmMetrics["MSpanInuse"] = float64(rtm.MSpanInuse)
	rtmMetrics["MSpanSys"] = float64(rtm.MSpanSys)
	rtmMetrics["Mallocs"] = float64(rtm.Mallocs)
	rtmMetrics["NextGC"] = float64(rtm.NextGC)
	rtmMetrics["NumForcedGC"] = float64(rtm.NumForcedGC)
	rtmMetrics["NumGC"] = float64(rtm.NumGC)
	rtmMetrics["OtherSys"] = float64(rtm.OtherSys)
	rtmMetrics["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	rtmMetrics["StackInuse"] = float64(rtm.StackInuse)
	rtmMetrics["StackSys"] = float64(rtm.StackSys)
	rtmMetrics["Sys"] = float64(rtm.Sys)
	rtmMetrics["TotalAlloc"] = float64(rtm.TotalAlloc)

	// Увеличиваем счетчик
	atomic.AddInt64(&a.pollCount, 1)

	// runtime метрики
	var metricsToSend []common.Metrics

	for metricName, metricValue := range rtmMetrics {
		m := common.Metrics{
			ID:    metricName,
			MType: string(common.GaugeMetric),
			Value: &metricValue,
		}

		metricsToSend = append(metricsToSend, m)
	}

	a.workChan <- metricsToSend

	// Сохраняем кол-во считываний
	pollCount := a.pollCount
	pollCountMetric := common.Metrics{
		ID:    "PollCount",
		MType: string(common.CounterMetric),
		Delta: &pollCount,
	}

	a.workChan <- []common.Metrics{pollCountMetric}

	// Сохраняем случайное значение
	randomValue := float64(rand.Intn(10000))
	randomValueMetric := common.Metrics{
		ID:    "RandomValue",
		MType: string(common.GaugeMetric),
		Value: &randomValue,
	}

	a.workChan <- []common.Metrics{randomValueMetric}
}

// Сбор метрик из пакета gopsutil
func (a *MetricAgent) PollPS() {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Сохраняем всю доступную память
	totalMemory := float64(v.Total)
	totalMemoryMetric := common.Metrics{
		ID:    "TotalMemory",
		MType: string(common.GaugeMetric),
		Value: &totalMemory,
	}

	a.workChan <- []common.Metrics{totalMemoryMetric}

	// Сохраняем свободную память
	freeMemory := float64(v.Free)
	freeMemoryMetric := common.Metrics{
		ID:    "FreeMemory",
		MType: string(common.GaugeMetric),
		Value: &freeMemory,
	}

	a.workChan <- []common.Metrics{freeMemoryMetric}

	// Сохраняем статистику по процессору
	var cpuMetrics []common.Metrics

	cpuStats, err := cpu.Percent(0, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	for cpuNum, cpuPercent := range cpuStats {
		m := common.Metrics{
			ID:    fmt.Sprintf("CPUutilization%d", cpuNum+1),
			MType: string(common.GaugeMetric),
			Value: &cpuPercent,
		}
		cpuMetrics = append(cpuMetrics, m)
	}

	a.workChan <- cpuMetrics
}
