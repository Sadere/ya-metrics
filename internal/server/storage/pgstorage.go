package storage

import (
	"context"

	"github.com/Sadere/ya-metrics/internal/common"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Хранение данных метрик в СУБД PostgreSQL
type PgMetricRepository struct {
	db *sqlx.DB
}

func NewPgRepository(db *sqlx.DB) *PgMetricRepository {
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

func (m PgMetricRepository) Set(metric common.Metrics) error {
	updateResult, err := m.db.Exec("UPDATE metrics SET delta = $1, value = $2 WHERE name = $3",
		metric.Delta,
		metric.Value,
		metric.ID,
	)
	if err != nil {
		return err
	}

	// Проверяем обновление данных метрики
	rowsAffected, err := updateResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		return nil
	}

	// Метрики не существует в БД, добавляем
	_, err = m.db.Exec(
		"INSERT INTO metrics (name, mtype, delta, value) VALUES ($1, $2, $3, $4)",
		metric.ID,
		metric.MType,
		metric.Delta,
		metric.Value,
	)

	return err
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
	ctx := context.Background()
	defer ctx.Done()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO metrics (name, mtype, delta, value) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		updateResult, err := tx.ExecContext(ctx, "UPDATE metrics SET delta = $1, value = $2 WHERE name = $3",
			metric.Delta,
			metric.Value,
			metric.ID,
		)
		if err != nil {
			return tx.Rollback()
		}

		rowsAffected, err := updateResult.RowsAffected()
		if err != nil {
			return tx.Rollback()
		}

		// Если успешно обновили идем к следующей метрике
		if rowsAffected > 0 {
			continue
		}

		// Добавляем метрику
		_, err = stmt.ExecContext(ctx,
			metric.ID,
			metric.MType,
			metric.Delta,
			metric.Value,
		)

		if err != nil {
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
