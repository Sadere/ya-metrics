// Пакет config считывает настройки агента из коммандной строки и переменных окружения, создает
// структуру Config со всеми настройками
package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Sadere/ya-metrics/internal/common"
)

// Параметры по умолчанию
const (
	DefaultPollInterval   = 2  // Интервал запроса метрик по умолчанию
	DefaultReportInterval = 10 // Интервал отправки метрики по умолчанию
	DefaultRateLimit      = 5  // Ограничение отправки по умолчанию
)

// Хранит настройки агента
type Config struct {
	ServerAddress common.NetAddress // Адрес сервера для отправки метрик

	PollInterval, // Интервал запроса метрик системы в секундах
	ReportInterval int // Интервал отправки метрик на сервер в секундах
	HashKey        string // Ключ для хеширования тела запроса
	RateLimit      int    // Кол-во одновременных отправок на сервер (кол-во воркеров)
	PubKeyFilePath string // Путь к файлу публичного ключа шифрования в формате PEM
}

// Возвращает структура конфига с установленными настройками
func NewConfig() Config {
	newConfig := Config{
		ServerAddress: common.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}

	// Парсим аргументы командной строки

	flag.IntVar(&newConfig.PollInterval, "p", DefaultPollInterval, "Частота сбора метрик")
	flag.IntVar(&newConfig.ReportInterval, "r", DefaultReportInterval, "Частота опроса сервера в секундах")
	flag.Var(&newConfig.ServerAddress, "a", "Адрес сервера")
	flag.StringVar(&newConfig.HashKey, "k", "", "Ключ для хеширования передаваемых данных")
	flag.IntVar(&newConfig.RateLimit, "l", DefaultRateLimit, "Лимит одновременных отправок на сервер")
	flag.StringVar(&newConfig.PubKeyFilePath, "crypto-key", "", "Путь к файлу публичного ключа шифрования в формате PEM")
	flag.Parse()

	// Берем опции из переменных окружения

	if envAddr := os.Getenv("ADDRESS"); len(envAddr) > 0 {
		err := newConfig.ServerAddress.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, ADDRESS = %s", envAddr)
		}
	}

	if envPollInt := os.Getenv("POLL_INTERVAL"); len(envPollInt) > 0 {
		number, err := strconv.Atoi(envPollInt)
		if err != nil {
			number = DefaultReportInterval
		}
		newConfig.PollInterval = number
	}

	if envReportInt := os.Getenv("REPORT_INTERVAL"); len(envReportInt) > 0 {
		number, err := strconv.Atoi(envReportInt)
		if err != nil {
			number = DefaultReportInterval
		}
		newConfig.ReportInterval = number
	}

	if envKey := os.Getenv("KEY"); len(envKey) > 0 {
		newConfig.HashKey = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); len(envRateLimit) > 0 {
		number, err := strconv.Atoi(envRateLimit)
		if err != nil {
			number = DefaultRateLimit
		}
		newConfig.RateLimit = number
	}

	if envPubKey := os.Getenv("CRYPTO_KEY"); len(envPubKey) > 0 {
		newConfig.PubKeyFilePath = envPubKey
	}

	return newConfig
}
