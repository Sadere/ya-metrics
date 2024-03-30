package config

import (
	"flag"
	"log"
	"os"

	"github.com/Sadere/ya-metrics/internal/common"
)

type Config struct {
	Address common.NetAddress
	LogLevel string
}

func NewConfig() Config {
	newConfig := Config{
		Address: common.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}

	flag.StringVar(&newConfig.LogLevel, "v", "fatal", "Уровень лога, возможные значения: debug, info, warn, error, dpanic, panic, fatal")
	flag.Var(&newConfig.Address, "a", "Адрес сервера")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); len(envAddr) > 0 {
		err := newConfig.Address.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, ADDRESS = %s", envAddr)
		}
	}

	return newConfig
}