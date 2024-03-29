package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Sadere/ya-metrics/internal/server/config"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     config.Config
	repository storage.MetricRepository
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

	// Загружаем HTML шаблоны
	execFile, _ := os.Executable()
	execPath := filepath.Dir(execFile)
	r.LoadHTMLGlob(execPath + "/../../templates/*")

	return r.Run(s.config.Address.String())
}

func Run() {
	server := &Server{}
	server.config = config.NewConfig()
	server.repository = storage.NewMemRepository()

	err := server.StartServer()
	if err != nil {
		log.Fatalf("couldn't launch server: %s", err.Error())
	}
}
