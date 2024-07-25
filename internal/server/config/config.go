// Пакет config считывает настройки для сервера
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

const (
	DefaultStoreInterval   = 300               // Значение по умолчанию для интервала записи
	DefaultFileStoragePath = "metrics-db.json" // Файл для хранения данных метрик по умолчанию
)

// Хранит настройки сервера
type Config struct {
	Address         common.NetAddress `json:"address"` // Адрес сервера
	LogLevel        string            // Уровень логирования
	StoreInterval   int               `json:"store_interval"` // Интервал в секундах через сколько сервер должен сохранять состояние в файл
	FileStoragePath string            `json:"store_file"`     // Путь к файлу
	Restore         bool              `json:"restore"`        // Восстанавливать данные из файла
	PostgresDSN     string            `json:"database_dsn"`   // DSN строка для подключения к бд
	HashKey         string            // Ключ для проверки хеша и хеширования ответов сервера
	PrivateKeyPath  string            `json:"crypto_key"` // Путь к файлу приватного ключа в формате PEM
}

func NewConfig() Config {
	newConfig := Config{
		Address: common.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		StoreInterval:   DefaultStoreInterval,
		FileStoragePath: DefaultFileStoragePath,
		Restore:         true,
	}

	var (
		flagAddress         common.NetAddress
		flagStoreInterval   int
		flagFileStoragePath string
		flagRestore         bool
		flagPostgresDSN     string
		flagHashKey         string
		flagPrivateKeyPath  string

		cfgFilePath string
	)

	flag.StringVar(&newConfig.LogLevel, "v", "fatal", "Уровень лога, возможные значения: debug, info, warn, error, dpanic, panic, fatal")
	flag.Var(&flagAddress, "a", "Адрес сервера")
	flag.IntVar(&flagStoreInterval, "i", 0, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск (значение 0 делает запись синхронной)")
	flag.StringVar(&flagFileStoragePath, "f", "", "Путь к файлу, хранящему данные метрик")
	flag.BoolVar(&flagRestore, "r", false, "Флаг, указывающий нужно ли восстанавливать данные из файла")
	flag.StringVar(&flagPostgresDSN, "d", "", "DSN для postgresql")
	flag.StringVar(&flagHashKey, "k", "", "Ключ для проверки хеша и хеширования ответов сервера")
	flag.StringVar(&flagPrivateKeyPath, "crypto-key", "", "Путь к файлу приватного ключа в формате PEM")
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

	// Конфиг из переменных окружений

	if envAddr := os.Getenv("ADDRESS"); len(envAddr) > 0 {
		err := newConfig.Address.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, ADDRESS = %s", envAddr)
		}
	} else if len(flagAddress.Host) > 0 {
		newConfig.Address = flagAddress
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); len(envStoreInterval) > 0 {
		number, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			number = DefaultStoreInterval
		}
		newConfig.StoreInterval = number
	} else if flagStoreInterval > 0 {
		newConfig.StoreInterval = flagStoreInterval
	}

	if envFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		newConfig.FileStoragePath = envFilePath
	} else if len(flagFileStoragePath) > 0 {
		newConfig.FileStoragePath = flagFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); len(envRestore) > 0 {
		newConfig.Restore = envRestore == "true"
	} else if flagRestore {
		newConfig.Restore = true
	}

	if envDSN := os.Getenv("DATABASE_DSN"); len(envDSN) > 0 {
		newConfig.PostgresDSN = envDSN
	} else if len(flagPostgresDSN) > 0 {
		newConfig.PostgresDSN = flagPostgresDSN
	}

	if envKey := os.Getenv("KEY"); len(envKey) > 0 {
		newConfig.HashKey = envKey
	} else if len(flagHashKey) > 0 {
		newConfig.HashKey = flagHashKey
	}

	if envPrivateKey := os.Getenv("CRYPTO_KEY"); len(envPrivateKey) > 0 {
		newConfig.PrivateKeyPath = envPrivateKey
	} else if len(flagPrivateKeyPath) > 0 {
		newConfig.PrivateKeyPath = flagPrivateKeyPath
	}

	fmt.Printf("address = %s\n", newConfig.Address.String())
	fmt.Printf("log level = %s\n", newConfig.LogLevel)
	fmt.Printf("store interval = %d sec\n", newConfig.StoreInterval)
	fmt.Printf("file storage path = %s\n", newConfig.FileStoragePath)
	fmt.Printf("private key path = %s\n", newConfig.PrivateKeyPath)

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
