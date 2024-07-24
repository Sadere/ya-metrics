package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/config"
	"github.com/Sadere/ya-metrics/internal/server/logger"
	"github.com/Sadere/ya-metrics/internal/server/middleware"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Информация о сборке
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// Основная структура сервера
type Server struct {
	config      config.Config            // Конфиг сервера
	repository  storage.MetricRepository // Репозиторий метрики
	fileManager *storage.FileManager     // Менеджер файла хранящий метрики
	log         *zap.Logger              // Лог
	db          *sqlx.DB                 // Указатель на соединение с БД
}

func (s *Server) setupRouter() *gin.Engine {
	r := gin.New()

	// Подключаем логи
	r.Use(middleware.Logger(s.log))

	// Стандартный обработчик паники
	r.Use(gin.Recovery())

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

	return r
}

func (s *Server) restoreState() {
	restoredState, err := s.fileManager.ReadMetrics()
	if err != nil {
		s.log.Sugar().Errorf("unable to read state from file: %s", err.Error())
	}

	metricsData := make(map[string]common.Metrics)

	for _, m := range restoredState {
		metricsData[m.ID] = m
	}

	err = s.repository.SetData(metricsData)
	if err != nil {
		s.log.Sugar().Errorf("unable to restore state: %s", err.Error())
	}
}

func (s *Server) saveState() {
	metrics, err := s.repository.GetData()
	if err != nil {
		s.log.Sugar().Errorf("unable to read state for saving: %s", err.Error())
	}

	savedState := make([]common.Metrics, 0)

	for _, m := range metrics {
		savedState = append(savedState, m)
	}

	if err := s.fileManager.WriteMetrics(savedState); err != nil {
		s.log.Sugar().Errorf("unable to save state: %s", err.Error())
	}
}

// Запуск http-сервера
func (s *Server) StartServer() {
	// Инициализируем роутер
	r := s.setupRouter()

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

	// Ловим сигналы отключения сервера
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Sugar().Infoln("shutdown server ...")

	// Сохраняем состояние сервера перед выходом
	if s.config.StoreInterval == 0 {
		s.saveState()
	}
}

// Инициализация логов
func (s *Server) InitLogging() {
	zapLogger, err := logger.NewZapLogger(s.config.LogLevel)
	if err != nil {
		log.Fatal("Couldn't initialize zap logger")
	}

	s.log = zapLogger
}

// Точка входа в сервер, настройка и запуск сервера
func Run() {
	// Выводим информацию о сборке
	fmt.Print(common.BuildInfo(buildVersion, buildDate, buildCommit))

	server := &Server{}
	server.config = config.NewConfig()
	server.fileManager = storage.NewFileManager(server.config.FileStoragePath)

	// Выбираем хранилище
	if len(server.config.PostgresDSN) > 0 {
		db, err := sqlx.Connect("pgx", server.config.PostgresDSN)
		if err != nil {
			log.Fatal(err.Error())
		}

		server.db = db
		server.repository = storage.NewPgRepository(db)
	} else {
		server.repository = storage.NewMemRepository()
	}

	// Инициализируем логи
	server.InitLogging()

	// Сохранение/восстановление состояния из файла
	if len(server.config.FileStoragePath) > 0 {
		// Восстанавливаем данные из файла
		if server.config.Restore {
			server.restoreState()
		}

		// Сохраняем состояние сервера в файле, если в конфиге указано интервальное сохранение
		if server.config.StoreInterval > 0 {
			go func() {
				for {
					time.Sleep(time.Second * time.Duration(server.config.StoreInterval))

					server.saveState()
				}
			}()
		}
	}

	// Запускаем сервер
	server.StartServer()
}
