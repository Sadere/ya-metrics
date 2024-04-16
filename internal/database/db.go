package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const MAX_RETRIES = 3

var (
	ErrDbConnection = errors.New("")
)

func TryQueryRow(db *sqlx.DB, sql string, args ...any) (*sql.Row, error) {
	var err error
	timeOut := 1

	for tryCount := 0; tryCount < MAX_RETRIES; tryCount++ {
		row := db.QueryRow(sql, args...)
		err = row.Err()

		if err == nil {
			return row, nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && !pgerrcode.IsConnectionException(pgErr.Code) {
			return nil, err
		}

		time.Sleep(time.Duration(timeOut) * time.Second)
		timeOut += 2
	}

	return nil, ErrDbConnection
}


func TryExec(db *sqlx.DB, sql string, args ...any) (sql.Result, error) {	
	timeOut := 1

	for tryCount := 0; tryCount < MAX_RETRIES; tryCount++ {
		result, err := db.Exec(sql, args...)
		if err == nil {
			return result, nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && !pgerrcode.IsConnectionException(pgErr.Code) {
			return nil, err
		}

		time.Sleep(time.Duration(timeOut) * time.Second)
		timeOut += 2
	}

	return nil, ErrDbConnection
}
