package storage

import (
	"errors"
	"fmt"
	"strconv"
)

type Storage interface {
	Get(string) (string, error)
	Set(string, string) error

	GetInt64(string) (int64, error)
	SetInt64(string, int64) error

	GetFloat64(string) (float64, error)
	SetFlaot64(string, float64) error
}

type MemStorage struct {
	data map[string]string
}

func New() *MemStorage {
	return &MemStorage{data: make(map[string]string)}
}

func (m MemStorage) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key shouldn't be empty")
	}

	value, ok := m.data[key]

	if !ok {
		return "", fmt.Errorf("no data with %s key", key)
	}

	return value, nil
}

func (m MemStorage) Set(key string, value string) error {
	if len(key) == 0 {
		return errors.New("key shouldn't be empty")
	}

	m.data[key] = value

	return nil
}

func (m MemStorage) GetInt64(key string) (int64, error) {
	value, err := m.Get(key);
	if err != nil {
		return 0, err
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (m MemStorage) SetInt64(key string, value int64) error {
	setValue := strconv.FormatInt(value, 10)

	return m.Set(key, setValue)
}

func (m MemStorage) GetFloat64(key string) (float64, error) {
	value, err := m.Get(key);
	if err != nil {
		return 0, err
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (m MemStorage) SetFloat64(key string, value float64) error {
	setValue := strconv.FormatFloat(value, 'g', -1, 64)

	return m.Set(key, setValue)
}