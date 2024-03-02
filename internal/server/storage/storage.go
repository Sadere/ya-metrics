package storage

// Интерфейс для хранения данных о метриках
type Storage interface {
	Get(string) (string, error)
	Set(string, string) error

	GetInt64(string) (int64, error)
	SetInt64(string, int64) error

	GetFloat64(string) (float64, error)
	SetFloat64(string, float64) error
}