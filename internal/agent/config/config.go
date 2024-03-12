package config

import (
	"os"
	"flag"
	"strconv"

	"github.com/Sadere/ya-metrics/internal/server"
)

type Config struct {
	Host string
	Port int

	PollInterval,
	ReportInterval int
}

func NewConfig() Config {
	var newConfig Config

	// Конфигурируем адрес сервера
	defaultHost := "localhost"
	defaultPort := 8080

	addr := &server.NetAddress{}
	addr.Host = defaultHost
	addr.Port = defaultPort

	envAddr, hasEnvAddr := os.LookupEnv("ADDRESS")

	if hasEnvAddr {
		err := addr.Set(envAddr)
		if err != nil {
			addr.Host = defaultHost
			addr.Port = defaultPort
		}
	} else {
		flag.Var(addr, "a", "Адрес сервера")
	}

	// Конфигурируем задержки
	flagReportInterval := flag.Int("r", 10, "Частота опроса сервера в секундах")
	flagPollInterval := flag.Int("p", 2, "Частота сбора метрик")
	flag.Parse()

	var optPollInterval, optReportInterval int

	// Частота опроса сервера
	envReportInterval, hasEnvReportInterval := os.LookupEnv("REPORT_INTERVAL")
	if hasEnvReportInterval {
		envInt, err := strconv.Atoi(envReportInterval)
		if err == nil {
			optReportInterval = envInt
		}
	} else {
		optReportInterval = *flagReportInterval
	}

	// Частота сбора метрик
	envPollInterval, hasEnvPollInterval := os.LookupEnv("POLL_INTERVAL")
	if hasEnvPollInterval {
		envInt, err := strconv.Atoi(envPollInterval)
		if err == nil {
			optPollInterval = envInt
		}
	} else {
		optPollInterval = *flagPollInterval
	}

	newConfig = Config{
		Host: addr.Host,
		Port: addr.Port,
		PollInterval: optPollInterval,
		ReportInterval: optReportInterval,
	}

	return newConfig
}