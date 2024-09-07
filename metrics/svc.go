package metrics

import (
	"context"
	"time"

	go_metric "github.com/hashicorp/go-metrics"
	go_metric_prometheus "github.com/hashicorp/go-metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ragpanda/go-toolkit/log"
)

type MetricsHub struct {
	args   *MetricsHubConfig
	config *go_metric.Config

	register *prometheus.Registry

	prometheusRegistry prometheus.Registerer
	prometheusSink     *go_metric_prometheus.PrometheusSink
	memSink            *go_metric.InmemSink
}

func NewMetricsHub(ctx context.Context, args *MetricsHubConfig) (*MetricsHub, error) {
	m := &MetricsHub{
		args: args,
	}
	err := m.Init()
	if err != nil {
		log.Error(ctx, "init metrics hub failed", err)
		return nil, err
	}
	return m, nil
}

type MetricsBackendType string

const (
	PrometheusBackendType MetricsBackendType = "prometheus"
	InmemBackendType      MetricsBackendType = "inmem"
)

// init it
func (self *MetricsHub) Init() error {
	if self.args == nil {
		self.args = &MetricsHubConfig{}
	}
	if self.args.ExpirationSec == 0 {
		self.args.ExpirationSec = 60 * 60 * 1
	}

	self.config = go_metric.DefaultConfig(self.args.ServiceName)

	var err error
	switch self.args.BackendType {
	case PrometheusBackendType:
		err = self.setPrometheusBackend()
	default:
		err = self.setInMemeBackend()
	}
	if err != nil {
		return err
	}

	return nil
}

// release it
func (self *MetricsHub) Release() {

	switch self.args.BackendType {
	case PrometheusBackendType:
		self.releasePrometheusBackend()
	default:
		self.releaseInMemeBackend()
	}

}

func (self *MetricsHub) setInMemeBackend() error {
	sink := go_metric.NewInmemSink(10*time.Second, 5*time.Minute)
	self.memSink = sink
	go_metric.NewGlobal(self.config, self.memSink)

	return nil
}

func (self *MetricsHub) setPrometheusBackend() error {
	reg := prometheus.NewRegistry()
	sink, err := go_metric_prometheus.NewPrometheusSinkFrom(go_metric_prometheus.PrometheusOpts{
		Expiration:         time.Duration(self.args.ExpirationSec) * time.Second,
		Registerer:         reg,
		GaugeDefinitions:   []go_metric_prometheus.GaugeDefinition{},
		SummaryDefinitions: []go_metric_prometheus.SummaryDefinition{},
		CounterDefinitions: []go_metric_prometheus.CounterDefinition{},
		Name:               self.args.ServiceName,
	})
	if err != nil {
		return err
	}
	self.prometheusRegistry = reg
	self.prometheusSink = sink
	self.register = prometheus.NewRegistry()

	go_metric.NewGlobal(self.config, self.memSink)
	return nil
}

func (self *MetricsHub) releaseInMemeBackend() {
	return
}

func (self *MetricsHub) releasePrometheusBackend() {
	self.register.Unregister(self.prometheusSink)
}
