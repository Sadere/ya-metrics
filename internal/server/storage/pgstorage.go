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
	var metric common.Metrics

	row := m.db.QueryRow("SELECT name, mtype, delta, value FROM metrics WHERE name = $1 AND mtype = $2", key, metricType)

	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		return metric, err
	}

	return metric, nil
}

func (m PgMetricRepository) Set(key string, metric common.Metrics) error {
	_, err := m.db.Exec(
		"INSERT INTO metrics (name, mtype, delta, value) VALUES ($1,$2,$3,$4)",
		key,
		metric.MType,
		metric.Delta,
		metric.Value,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m PgMetricRepository) GetData() (map[string]common.Metrics, error) {
	metrics := make(map[string]common.Metrics)

	rows, err := m.db.Query("SELECT name, mtype, delta, value FROM metrics")
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var m common.Metrics

		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err == nil {
			metrics[m.ID] = m
		}
	}

	if err = rows.Err(); err != nil {
		return metrics, nil
	}

	return metrics, nil
}

func (m PgMetricRepository) SetData(metrics map[string]common.Metrics) error {
	for key, metric := range metrics {
		if err := m.Set(key, metric); err != nil {
			return err
		}
	}

	return nil
}
