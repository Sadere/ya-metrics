package server

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	storage storage.Storage
}

type NetAddress struct {
	Host string
	Port int
}

func (addr *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

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

func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()

	// Обработка обновления метрик
	r.POST(`/update/:type/:metric/:value`, s.updateHandle)

	// Вывод метрики
	r.GET(`/value/:type/:metric`, s.getMetricHandle)

	// Вывод всех метрик в HTML
	r.GET(`/`, s.getAllMetricsHandle)

	return r
}

func (s *Server) StartServer() error {
	r := s.setupRouter()

	// Конфигурируем
	defaultHost := "localhost"
	defaultPort := 8080

	addr := &NetAddress{}
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
	
	flag.Parse()

	fmt.Println(addr)

	// Загружаем HTML шаблоны
	execFile, _ := os.Executable()
	execPath := filepath.Dir(execFile)
	r.LoadHTMLGlob(execPath + "/../../templates/*")

	return r.Run(addr.String())
}

func Run() {
	server := &Server{}
	server.storage = storage.NewMemStorage()

	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
