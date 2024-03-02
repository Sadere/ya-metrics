package server

import (
	"net/http"

	"github.com/Sadere/ya-metrics/internal/server/storage"
)

type Server struct {
	storage storage.Storage
}

// Middleware отсекающий все запросы кроме POST
func (s *Server) postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST method allowed", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func (s *Server) StartServer() error {
	mux := http.NewServeMux()

	// Обработка обновления метрик
	mux.Handle(`/update/gauge/`, s.postMiddleware(http.HandlerFunc(s.updateGaugeHandle)))
	mux.Handle(`/update/counter/`, s.postMiddleware(http.HandlerFunc(s.updateCounterHandle)))

	// Обработка остальных запросов
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	return http.ListenAndServe(`:8080`, mux)
}

func Run() {
	server := &Server{}
	server.storage = storage.NewMemStorage()

	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
