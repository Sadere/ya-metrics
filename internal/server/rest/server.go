package rest

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sadere/ya-metrics/internal/server/config"
	"github.com/Sadere/ya-metrics/internal/server/logger"
	"github.com/Sadere/ya-metrics/internal/server/middleware"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Основная структура сервера
type Server struct {
	config        config.Config          // Конфиг сервера
	metricService *service.MetricService // Сервис метрик
	log           *zap.Logger            // Лог
	db            *sqlx.DB               // Указатель на соединение с БД
}

func NewServer(
	cfg config.Config,
	mServ *service.MetricService,
	log *zap.Logger,
	db *sqlx.DB,
	) *Server {

	return &Server{
		config: cfg,
		metricService: mServ,
		log: log,
		db: db,
	}
}

func (s *Server) Start() error {
	// Инициализируем роутер
	r, err := s.setupRouter()

	if err != nil {
		return err
	}

	// Загружаем HTML шаблоны
	execFile, _ := os.Executable()
	execPath := filepath.Dir(execFile)
	r.LoadHTMLGlob(execPath + "/../../templates/*")

	srv := &http.Server{
		Addr:    s.config.Address.String(),
		Handler: r,
	}

	// Запускаем сервер в фоне
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Sugar().Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}

func (s *Server) setupRouter() (*gin.Engine, error) {
	r := gin.New()

	// Подключаем логи
	r.Use(middleware.Logger(s.log))

	// Стандартный обработчик паники
	r.Use(gin.Recovery())

	if len(s.config.TrustedSubnet) > 0 {
		IPMiddleware, err := middleware.ValidateIP(s.config.TrustedSubnet)
		if err != nil {
			s.log.Sugar().Fatalf("failed to setup trusted subnet: %s", err)
			return nil, err
		}

		r.Use(IPMiddleware)
	}

	// Распаковываем запрос
	r.Use(middleware.GzipDecompress())

	// Проверка хеша
	r.Use(middleware.ValidateHash(s.config.HashKey))

	// Дешифровка запроса
	if len(s.config.PrivateKeyPath) > 0 {
		RSAMiddleware, err := middleware.Decrypt(s.config.PrivateKeyPath)
		if err != nil {
			s.log.Sugar().Errorf("unable to initialize RSA decryption middleware: %s", err.Error())
		} else {
			r.Use(RSAMiddleware)
		}
	}

	// Хеш ответа
	r.Use(middleware.HashResponse(s.config.HashKey))

	// Упаковываем ответ
	r.Use(middleware.GzipCompress())

	// Обработка обновления метрик
	r.POST(`/update/:type/:metric/:value`, s.updateHandle)
	r.POST(`/update/`, middleware.JSON(), s.updateHandleJSON)

	// Вывод метрики
	r.GET(`/value/:type/:metric`, s.getMetricHandle)
	r.POST(`/value/`, middleware.JSON(), s.getMetricHandleJSON)

	// Проверка подключения к бд
	r.GET(`/ping`, s.pingHandle)

	// Добавление нескольких метрик
	r.POST(`/updates/`, middleware.JSON(), s.updateBatchHandleJSON)

	// Вывод всех метрик в HTML
	r.GET(`/`, s.getAllMetricsHandle)

	return r, nil
}

// Инициализация логов
func (s *Server) InitLogging() {
	zapLogger, err := logger.NewZapLogger(s.config.LogLevel)
	if err != nil {
		log.Fatal("Couldn't initialize zap logger")
	}

	s.log = zapLogger
}