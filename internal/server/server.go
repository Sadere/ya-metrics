package server

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sadere/ya-metrics/internal/server/storage"
)

type Server struct {
	storage *storage.MemStorage
}

func urlParams(url *url.URL) []string {
	return strings.Split(url.Path, "/")[3:]
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

func (s *Server) updateGaugeHandle(res http.ResponseWriter, req *http.Request) {
	params := urlParams(req.URL)

	if len(params) < 2 {
		http.Error(res, "Insufficient parameters", http.StatusNotFound)
		return
	}

	name := params[0]
	value := params[1]

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.SetFloat64(name, valueFloat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) updateCounterHandle(res http.ResponseWriter, req *http.Request) {
	params := urlParams(req.URL)

	if len(params) < 2 {
		http.Error(res, "Insufficient parameters", http.StatusNotFound)
		return
	}

	name := params[0]

	oldValue, err := s.storage.GetInt64(name)
	if err != nil {
		oldValue = 0
	}

	addValue, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.SetInt64(name, addValue+oldValue)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	mux.Handle(`/update/gauge/`, s.postMiddleware(http.HandlerFunc(s.updateGaugeHandle)))
	mux.Handle(`/update/counter/`, s.postMiddleware(http.HandlerFunc(s.updateCounterHandle)))

	return http.ListenAndServe(`:8080`, mux)
}

func Run() {
	server := &Server{}
	server.storage = storage.New()

	err := server.RunServer()
	if err != nil {
		panic(err)
	}
}
