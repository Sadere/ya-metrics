// Пакет config считывает настройки агента из коммандной строки и переменных окружения, создает
// структуру Config со всеми настройками
package config

import (
	"encoding/json"
	"flag"
	"fmt"
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
	ServerAddress  common.NetAddress `json:"address"`         // Адрес сервера для отправки метрик
	PollInterval   int               `json:"poll_interval"`   // Интервал запроса метрик системы в секундах
	ReportInterval int               `json:"report_interval"` // Интервал отправки метрик на сервер в секундах
	HashKey        string            // Ключ для хеширования тела запроса
	RateLimit      int               // Кол-во одновременных отправок на сервер (кол-во воркеров)
	PubKeyFilePath string            `json:"crypto_key"` // Путь к файлу публичного ключа шифрования в формате PEM
	HostAddress    string
}

// Возвращает структура конфига с установленными настройками
func NewConfig() Config {
	newConfig := Config{
		ServerAddress: common.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		PollInterval:   DefaultPollInterval,
		ReportInterval: DefaultReportInterval,
		RateLimit:      DefaultRateLimit,
	}

	var (
		flagPollInterval   int
		flagReportInterval int
		flagServerAddress  common.NetAddress
		flagHashKey        string
		flagRateLimit      int
		flagPubKeyFilePath string

		cfgFilePath string
	)

	// Парсим аргументы командной строки

	flag.IntVar(&flagPollInterval, "p", 0, "Частота сбора метрик")
	flag.IntVar(&flagReportInterval, "r", 0, "Частота опроса сервера в секундах")
	flag.Var(&flagServerAddress, "a", "Адрес сервера")
	flag.StringVar(&flagHashKey, "k", "", "Ключ для хеширования передаваемых данных")
	flag.IntVar(&flagRateLimit, "l", 0, "Лимит одновременных отправок на сервер")
	flag.StringVar(&flagPubKeyFilePath, "crypto-key", "", "Путь к файлу публичного ключа шифрования в формате PEM")
	flag.StringVar(&cfgFilePath, "c", "", "Путь к файлу конфига")
	flag.Parse()

	// Берем конфигурацию из файла, если передан путь до конфига
	if envCfgFile := os.Getenv("CONFIG"); len(envCfgFile) > 0 {
		cfgFilePath = envCfgFile
	}

	if len(cfgFilePath) > 0 {
		cfgFromFile, err := FromFile(cfgFilePath)
		if err == nil {
			newConfig = cfgFromFile
		}
	}

	// Берем опции из переменных окружения в приоритете

	if envAddr := os.Getenv("ADDRESS"); len(envAddr) > 0 {
		err := newConfig.ServerAddress.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, ADDRESS = %s", envAddr)
		}
	} else if len(flagServerAddress.Host) > 0 {
		newConfig.ServerAddress = flagServerAddress
	}

	if envPollInt := os.Getenv("POLL_INTERVAL"); len(envPollInt) > 0 {
		number, err := strconv.Atoi(envPollInt)
		if err != nil {
			number = DefaultReportInterval
		}
		newConfig.PollInterval = number
	} else if flagPollInterval > 0 {
		newConfig.PollInterval = flagPollInterval
	}

	if envReportInt := os.Getenv("REPORT_INTERVAL"); len(envReportInt) > 0 {
		number, err := strconv.Atoi(envReportInt)
		if err != nil {
			number = DefaultReportInterval
		}
		newConfig.ReportInterval = number
	} else if flagReportInterval > 0 {
		newConfig.ReportInterval = flagReportInterval
	}

	if envKey := os.Getenv("KEY"); len(envKey) > 0 {
		newConfig.HashKey = envKey
	} else if len(flagHashKey) > 0 {
		newConfig.HashKey = flagHashKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); len(envRateLimit) > 0 {
		number, err := strconv.Atoi(envRateLimit)
		if err != nil {
			number = DefaultRateLimit
		}
		newConfig.RateLimit = number
	} else if flagRateLimit > 0 {
		newConfig.RateLimit = flagRateLimit
	}

	if envPubKey := os.Getenv("CRYPTO_KEY"); len(envPubKey) > 0 {
		newConfig.PubKeyFilePath = envPubKey
	} else if len(flagPubKeyFilePath) > 0 {
		newConfig.PubKeyFilePath = flagPubKeyFilePath
	}

	fmt.Printf("server address = %s\n", newConfig.ServerAddress.String())
	fmt.Printf("poll interval = %d sec\n", newConfig.PollInterval)
	fmt.Printf("report interval = %d sec\n", newConfig.ReportInterval)
	fmt.Printf("rate limit = %d\n", newConfig.RateLimit)
	fmt.Printf("public key path = %s\n", newConfig.PubKeyFilePath)

	return newConfig
}

// Получает конфиг из файла filePath
func FromFile(filePath string) (Config, error) {
	cfg := Config{}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
