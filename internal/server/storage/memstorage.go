package storage

import (
	"errors"
	"fmt"
)

// Хранение данных метрик в памяти
type MemStorage struct {
	data map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]interface{})}
}

func (m MemStorage) Get(key string) (interface{}, error) {
	if len(key) == 0 {
		return "", errors.New("key shouldn't be empty")
	}

	value, ok := m.data[key]

	if !ok {
		return "", fmt.Errorf("no data with %s key", key)
	}

	return value, nil
}

func (m MemStorage) Set(key string, value interface{}) error {
	if len(key) == 0 {
		return errors.New("key shouldn't be empty")
	}

	m.data[key] = value

	return nil
}

func (m MemStorage) GetData() map[string]interface{} {
	return m.data
}