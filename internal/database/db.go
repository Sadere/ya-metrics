package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

// Дефолтные параметры для backoff стратегии
const (
	InitialInterval = time.Second
	MaxRetries      = 3
)

// Ошибки БД
var (
	ErrDBConnection = errors.New("couldn't establish db connection") // Ошибка соединения к БД
)

func newBackoff() backoff.BackOff {
	return backoff.WithMaxRetries(
		backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(InitialInterval),
		),
		MaxRetries,
	)
}

// Обертка над QueryRow использующая backoff стратегию
func TryQueryRow(db *sqlx.DB, sql string, args ...any) (row *sql.Row, err error) {
	b := newBackoff()

	operation := func() error {
		row = db.QueryRow(sql, args...)

		return row.Err()
	}

	err = backoff.Retry(operation, b)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return nil, ErrDBConnection
	}

	return row, err
}

// Обертка над Exec использующая backoff стратегию
func TryExec(db *sqlx.DB, sql string, args ...any) (result sql.Result, err error) {
	b := newBackoff()

	operation := func() error {
		result, err = db.Exec(sql, args...)

		return err
	}

	err = backoff.Retry(operation, b)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return nil, ErrDBConnection
	}

	return result, err
}
