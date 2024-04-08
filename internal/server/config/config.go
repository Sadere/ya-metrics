package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Sadere/ya-metrics/internal/common"
)

const (
	DefaultStoreInterval   = 300 // Значение по умолчанию для интервала записи
	DefaultFileStoragePath = "metrics-db.json"
)

type Config struct {
	Address         common.NetAddress // Адрес сервера
	LogLevel        string            // Уровень логирования
	StoreInterval   int               // Интервал в секундах через сколько сервер должен сохранять состояние в файл
	FileStoragePath string            // Путь к файлу
	Restore         bool              // Восстанавливать данные из файла
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
	flag.IntVar(&newConfig.StoreInterval, "i", DefaultStoreInterval, "Bнтервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск (значение 0 делает запись синхронной)")
	flag.StringVar(&newConfig.FileStoragePath, "f", DefaultFileStoragePath, "Путь к файлу, хранящему данные метрик")
	flag.BoolVar(&newConfig.Restore, "r", true, "Флаг, указывающий нужно ли восстанавливать данные из файла")
	flag.Parse()

	// Конфиг из переменных окружений

	if envAddr := os.Getenv("ADDRESS"); len(envAddr) > 0 {
		err := newConfig.Address.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, ADDRESS = %s", envAddr)
		}
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); len(envStoreInterval) > 0 {
		number, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			number = DefaultStoreInterval
		}
		newConfig.StoreInterval = number
	}

	if envFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		newConfig.FileStoragePath = envFilePath
	}

	if envRestore := os.Getenv("RESTORE"); len(envRestore) > 0 {
		newConfig.Restore = envRestore == "true"
	}

	return newConfig
}
