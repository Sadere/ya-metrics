package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Example() {
	type metric struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}

	d := int64(100)
	m := metric{
		ID:    "exampleCounter",
		MType: "counter",
		Delta: &d,
	}
	postBody, _ := json.Marshal(m)

	// Пример отправки метрики
	_, err := http.Post("http://example.com:8080/update/", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		panic(err)
	}

	// Пример получения значения метрики
	var resultMetric metric

	resp, err := http.Post("http://example.com:8080/value/", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resultMetric)
	if err != nil {
		panic(err)
	}

	// resultMetric содержит структуру метрики

	// Пример отправки нескольких метрик
	var metrics []metric

	v := 123.456
	for i := 0; i < 10; i++ {
		metrics = append(metrics, metric{
			ID:    fmt.Sprintf("exampleGauge%d", i),
			MType: "gauge",
			Value: &v,
		})
	}

	postBody, _ = json.Marshal(metrics)

	_, err = http.Post("http://example.com:8080/updates/", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		panic(err)
	}

}
