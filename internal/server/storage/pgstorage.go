package storage

import (
	"database/sql"

	"github.com/Sadere/ya-metrics/internal/common"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Хранение данных метрик в СУБД PostgreSQL
type PgMetricRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgMetricRepository {
	return &PgMetricRepository{db: db}
}

func (m PgMetricRepository) Get(metricType common.MetricType, key string) (common.Metrics, error) {
	return common.Metrics{}, nil
}

func (m PgMetricRepository) Set(key string, metric common.Metrics) error {
	return nil
}

func (m PgMetricRepository) GetData() map[string]common.Metrics {
	return make(map[string]common.Metrics)
}

func (m PgMetricRepository) SetData(metrics map[string]common.Metrics) {
}
