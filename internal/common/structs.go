package common

import (
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (addr *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

func (addr *NetAddress) Set(flagValue string) error {
	addrParts := strings.Split(flagValue, ":")

	if len(addrParts) == 2 {
		addr.Host = addrParts[0]
		optPort, err := strconv.Atoi(addrParts[1])
		if err != nil {
			return err
		}

		addr.Port = optPort
	}

	return nil
}

type MetricType string

const (
	CounterMetric MetricType = "counter"
	GaugeMetric   MetricType = "gauge"
)


type Metric struct {
	Type MetricType
	ValueGauge float64
	ValueCounter int64
}