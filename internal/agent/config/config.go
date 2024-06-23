package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Sadere/ya-metrics/internal/common"
)

const (
	DefaultPollInterval   = 2
	DefaultReportInterval = 10
	DefaultRateLimit      = 5
)

type Config struct {
	ServerAddress common.NetAddress

	PollInterval,
	ReportInterval int
	CryptoKey string
	RateLimit int
}

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
	flag.StringVar(&newConfig.CryptoKey, "k", "", "Ключ для хеширования передаваемых данных")
	flag.IntVar(&newConfig.RateLimit, "l", DefaultRateLimit, "Лимит одновременных отправок на сервер")
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
		newConfig.CryptoKey = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); len(envRateLimit) > 0 {
		number, err := strconv.Atoi(envRateLimit)
		if err != nil {
			number = DefaultRateLimit
		}
		newConfig.RateLimit = number
	}

	return newConfig
}
