package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func urlParams(url *url.URL) []string {
	urlParts := strings.Split(url.Path, "/")

	if len(urlParts) > 3 {
		return urlParts[3:]
	}

	return []string{}
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

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Println(s.storage)
}
