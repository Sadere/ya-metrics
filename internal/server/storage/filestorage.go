package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/config"
	"github.com/Sadere/ya-metrics/internal/server/logger"
)

// Оболочка над MemRepository позволяющая хранить и читать данные из файла
type FileMetricRepository struct {
	mem           *MemMetricRepository
	mutex         *sync.Mutex
	StoreInterval time.Duration
	FilePath      string
	Restore       bool
}

func NewFileRepository(cfg config.Config) (*FileMetricRepository, error) {
	newRepository := &FileMetricRepository{
		mem:           NewMemRepository(),
		mutex:         &sync.Mutex{},
		StoreInterval: time.Second * time.Duration(cfg.StoreInterval),
		FilePath:      cfg.FileStoragePath,
		Restore:       cfg.Restore,
	}

	if cfg.Restore {
		err := newRepository.restoreMetrics()
		if err != nil {
			return nil, err
		}
	}

	// Запись в фоне по интервалу
	if newRepository.StoreInterval.Seconds() > 0 {
		go func() {
			time.Sleep(newRepository.StoreInterval)

			if err := newRepository.writeMetrics(); err != nil {
				logger.Log.Sugar().Fatalf("failed to save metrics: %s", err.Error())
			}
		}()
	}

	return newRepository, nil
}

func (f FileMetricRepository) Get(metricType common.MetricType, key string) (common.Metrics, error) {
	return f.mem.Get(metricType, key)
}

func (f FileMetricRepository) Set(key string, metric common.Metrics) error {
	err := f.mem.Set(key, metric)
	if err != nil {
		return err
	}

	// Синхронная запись
	if f.StoreInterval.Seconds() == 0 {
		err := f.writeMetrics()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f FileMetricRepository) GetData() map[string]common.Metrics {
	return f.mem.GetData()
}

func (f FileMetricRepository) restoreMetrics() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var metric common.Metrics

	file, err := os.OpenFile(f.FilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	for {
		err := decoder.Decode(&metric)
		if err != nil {
			break
		}

		err = f.Set(metric.ID, metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f FileMetricRepository) writeMetrics() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	for _, metric := range f.mem.GetData() {
		err := encoder.Encode(metric)
		if err != nil {
			return err
		}
	}

	return nil
}
