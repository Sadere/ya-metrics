package storage

// Интерфейс для хранения данных о метриках
type Storage interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error

	GetData() map[string]interface{}
}