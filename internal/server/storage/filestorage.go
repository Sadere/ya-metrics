package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Sadere/ya-metrics/internal/common"
)

// Структура для сохранения и чтения данных метрик в файле
type FileManager struct {
	mutex    *sync.Mutex
	filePath string
}

func NewFileManager(filePath string) *FileManager {
	return &FileManager{
		mutex:    &sync.Mutex{},
		filePath: filePath,
	}
}

// Возвращает массив метрик, прочитанный из файла
func (f FileManager) ReadMetrics() ([]common.Metrics, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	result := make([]common.Metrics, 0)

	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	for {
		var metric common.Metrics

		err := decoder.Decode(&metric)
		if err != nil {
			break
		}

		result = append(result, metric)
	}

	return result, nil
}

// Сохраняет массив метрик в файл
func (f FileManager) WriteMetrics(metrics []common.Metrics) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	for _, metric := range metrics {
		err := encoder.Encode(metric)
		if err != nil {
			return err
		}
	}

	return nil
}
