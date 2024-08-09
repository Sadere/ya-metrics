package storage

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func NewMock(t *testing.T) (*sqlx.DB, sqlxmock.Sqlmock) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to create mock db: %s", err)
	}

	return db, mock
}

func TestGet(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	mType := common.GaugeMetric
	key := "test_gauge"
	v := float64(10.20)

	expectedMetric := common.Metrics{
		ID:    key,
		MType: string(mType),
		Value: &v,
	}

	rows := sqlxmock.NewRows([]string{"name", "mtype", "delta", "value"}).AddRow(key, string(mType), nil, v)
	mock.ExpectQuery("SELECT (.+) FROM metrics WHERE").WithArgs(key, mType).WillReturnRows(rows)

	res, err := repo.Get(mType, key)

	assert.NoError(t, err)

	if !reflect.DeepEqual(expectedMetric, res) {
		t.Errorf("want metric = %v | got = %v", expectedMetric, res)
	}
}

func TestGetError(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	mType := common.CounterMetric
	key := "error_counter"

	rows := sqlxmock.NewRows([]string{"name", "mtype", "delta", "value"})
	mock.ExpectQuery("SELECT (.+) FROM metrics WHERE").WithArgs(key, mType).WillReturnRows(rows)

	_, err := repo.Get(mType, key)

	assert.Error(t, err)
}

// Проверяет Set функцию в случае когда обновляем существующую метрику
func TestSetUpdated(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	d := int64(100)
	metric := common.Metrics{
		ID:    "test_counter",
		MType: string(common.CounterMetric),
		Delta: &d,
	}

	updRes := sqlxmock.NewResult(0, 1)
	mock.ExpectExec("UPDATE metrics").WithArgs(metric.Delta, metric.Value, metric.ID).WillReturnResult(updRes)

	err := repo.Set(metric)

	assert.NoError(t, err)
}

// Проверяет Set функцию в случае когда добавляем новую метрику
func TestSetInsert(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	d := int64(100)
	metric := common.Metrics{
		ID:    "test_counter",
		MType: string(common.CounterMetric),
		Delta: &d,
	}

	updRes := sqlxmock.NewResult(0, 0)
	mock.ExpectExec("UPDATE metrics").WithArgs(metric.Delta, metric.Value, metric.ID).WillReturnResult(updRes)

	insRes := sqlxmock.NewResult(1, 1)
	mock.ExpectExec("INSERT INTO metrics").WithArgs(metric.ID, metric.MType, metric.Delta, metric.Value).WillReturnResult(insRes)

	err := repo.Set(metric)

	assert.NoError(t, err)
}

func TestGetData(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	v := float64(12.34)
	m := common.Metrics{
		ID:    "test_metric",
		MType: string(common.GaugeMetric),
		Value: &v,
	}

	expectedResult := map[string]common.Metrics{m.ID: m}

	rows := sqlxmock.NewRows([]string{"name", "mtype", "delta", "value"}).AddRow(m.ID, m.MType, nil, v)
	mock.ExpectQuery("SELECT (.+) FROM metrics").WillReturnRows(rows)

	res, err := repo.GetData()

	assert.NoError(t, err)

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("want result = %v | got = %v", expectedResult, res)
	}
}

func TestGetDataError(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM metrics").WillReturnError(errors.New("test error"))

	_, err := repo.GetData()

	assert.Error(t, err)
}

// Проверяет SetData функцию в случае когда обновляем существующую метрику
func TestSetDataUpdated(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	d := int64(100)
	metric := common.Metrics{
		ID:    "test_counter",
		MType: string(common.CounterMetric),
		Delta: &d,
	}

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO metrics")

	updRes := sqlxmock.NewResult(0, 1)
	mock.ExpectExec("UPDATE metrics").WithArgs(metric.Delta, metric.Value, metric.ID).WillReturnResult(updRes)

	mock.ExpectCommit()

	err := repo.SetData(map[string]common.Metrics{metric.ID: metric})

	assert.NoError(t, err)
}

// Проверяет SetData функцию в случае когда добавляем новую метрику
func TestSetDataInsert(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgRepository(db)

	d := int64(100)
	metric := common.Metrics{
		ID:    "test_counter",
		MType: string(common.CounterMetric),
		Delta: &d,
	}

	mock.ExpectBegin()

	insRes := sqlxmock.NewResult(1, 1)
	prep := mock.ExpectPrepare("INSERT INTO metrics")

	updRes := sqlxmock.NewResult(0, 0)
	mock.ExpectExec("UPDATE metrics").WithArgs(metric.Delta, metric.Value, metric.ID).WillReturnResult(updRes)

	prep.ExpectExec().WithArgs(metric.ID, metric.MType, metric.Delta, metric.Value).WillReturnResult(insRes)

	mock.ExpectCommit()

	err := repo.SetData(map[string]common.Metrics{metric.ID: metric})

	assert.NoError(t, err)
}
