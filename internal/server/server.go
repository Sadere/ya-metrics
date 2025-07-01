package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/config"
	"github.com/Sadere/ya-metrics/internal/server/grpc"
	"github.com/Sadere/ya-metrics/internal/server/logger"
	"github.com/Sadere/ya-metrics/internal/server/rest"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Информация о сборке
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

type MetricApp struct {
	config        config.Config
	log           *zap.Logger
	fileManager   *storage.FileManager   // Менеджер файла хранящий метрики
	metricService *service.MetricService // Сервис метрик
}

func (a *MetricApp) restoreState() {
	restoredState, err := a.fileManager.ReadMetrics()
	if err != nil {
		a.log.Sugar().Errorf("unable to read state from file: %s", err.Error())
	}

	metricsData := make(map[string]common.Metrics)

	for _, m := range restoredState {
		metricsData[m.ID] = m
	}

	err = a.metricService.SaveMetrics(metricsData)
	if err != nil {
		a.log.Sugar().Errorf("unable to restore state: %s", err.Error())
	}
}

func (a *MetricApp) saveState() {
	metrics, err := a.metricService.GetAllMetrics()
	if err != nil {
		a.log.Sugar().Errorf("unable to read state for saving: %s", err.Error())
	}

	savedState := make([]common.Metrics, 0)

	for _, m := range metrics {
		savedState = append(savedState, m)
	}

	if err := a.fileManager.WriteMetrics(savedState); err != nil {
		a.log.Sugar().Errorf("unable to save state: %s", err.Error())
	}
}

// Инициализация логов
func (a *MetricApp) InitLogging(logLevel string) {
	zapLogger, err := logger.NewZapLogger(logLevel)
	if err != nil {
		log.Fatal("Couldn't initialize zap logger")
	}

	a.log = zapLogger
}

// Точка входа в сервер, настройка и запуск сервера
func Run() {
	// Выводим информацию о сборке
	fmt.Print(common.BuildInfo(buildVersion, buildDate, buildCommit))

	cfg, err := config.NewConfig(os.Args[1:])
	if err != nil {
		log.Fatalf("failed to initialize config: %s", err)
	}

	var (
		rep storage.MetricRepository
		db  *sqlx.DB
	)

	// Выбираем хранилище
	if len(cfg.PostgresDSN) > 0 {
		db, err = sqlx.Connect("pgx", cfg.PostgresDSN)
		if err != nil {
			log.Fatal(err.Error())
		}

		rep = storage.NewPgRepository(db)
	} else {
		rep = storage.NewMemRepository()
	}

	metricService := service.NewMetricService(rep)

	app := &MetricApp{
		config:        cfg,
		metricService: metricService,
	}

	// Инициализируем логи
	app.InitLogging(cfg.LogLevel)

	app.fileManager = storage.NewFileManager(cfg.FileStoragePath)

	// Сохранение/восстановление состояния из файла
	if len(cfg.FileStoragePath) > 0 {
		// Восстанавливаем данные из файла
		if cfg.Restore {
			app.restoreState()
		}

		// Сохраняем состояние сервера в файле, если в конфиге указано интервальное сохранение
		if cfg.StoreInterval > 0 {
			go func() {
				for {
					time.Sleep(time.Second * time.Duration(cfg.StoreInterval))

					app.saveState()
				}
			}()
		}
	}

	if cfg.ServeGRPC {
		app.StartGRPC()
	} else {
		app.StartREST(db)
	}
}

func (a *MetricApp) StartGRPC() {
	listen, err := net.Listen("tcp", a.config.Address.String())
	if err != nil {
		log.Fatalln("failed to listen", err)
	}

	server := grpc.NewServer(a.config, a.metricService, a.log)

	// Регистрируем gRPC сервер
	srv, err := server.Register()
	if err != nil {
		a.log.Sugar().Fatalf("failed to register gRPC server: %s\n", err)
	}

	// Запускаем сервер в фоне
	go func() {
		if err := srv.Serve(listen); err != nil {
			a.log.Sugar().Fatalf("gRPC serve error: %s\n", err)
		}
	}()

	// Ловим сигналы отключения сервера
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	a.log.Sugar().Infoln("gRPC server shutdown ...")

	srv.GracefulStop()

	// Сохраняем состояние сервера перед выходом
	if a.config.StoreInterval == 0 {
		a.saveState()
	}
}

func (a *MetricApp) StartREST(db *sqlx.DB) {
	server := rest.NewServer(a.config, a.metricService, a.log, db)

	// Запускаем сервер
	err := server.Start()
	if err != nil {
		a.log.Sugar().Fatalf("couldn't start server: %s\n", err)
	}

	// Ловим сигналы отключения сервера
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	a.log.Sugar().Infoln("shutdown server ...")

	// Сохраняем состояние сервера перед выходом
	if a.config.StoreInterval == 0 {
		a.saveState()
	}
}
