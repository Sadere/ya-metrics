// Пакет хранит структуры, типы и константы общие для сервера и агента
package common

import (
	"fmt"
	"strconv"
	"strings"
)

// Адрес хоста в формате <host>:<port>
type NetAddress struct {
	Host string
	Port int
}

// Выводит адрес в виде строки <host>:<port>
func (addr *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

// Парсинг адреса из строки
func (addr *NetAddress) Set(flagValue string) error {
	addrParts := strings.Split(flagValue, ":")

	if len(addrParts) == 2 {
		addr.Host = addrParts[0]
		optPort, err := strconv.Atoi(addrParts[1])
		if err != nil {
			return err
		}

		addr.Port = optPort
	}

	return nil
}

// Адрес из JSON данных
func (addr *NetAddress) UnmarshalJSON(data []byte) error {
	value := string(data[1:len(data)-1])
	return addr.Set(value)
}

// Тип метрики
type MetricType string

const (
	CounterMetric MetricType = "counter"
	GaugeMetric   MetricType = "gauge"
	HashHeader               = "HashSHA256"
	AESKeyHeader             = "X-AES-Key"
)

// Структура для хранения одной метрики
type Metrics struct {
	ID    string   `json:"id" db:"name"`               // имя метрики
	MType string   `json:"type" db:"mtype"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
}
