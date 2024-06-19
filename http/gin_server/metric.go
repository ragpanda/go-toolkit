package gin_server

import "github.com/prometheus/client_golang/prometheus"

const HTTPMetricServerPrefix = "http.server."

var (
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPCurrentRequests *prometheus.GaugeVec
	HTTPRequestDuration *prometheus.HistogramVec
)

var (
	// 计数器：记录HTTP请求的总次数
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
	// 仪表：记录当前在处理的请求数量
	HttpSrvcurrentRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "current_requests",
			Help: "Number of requests currently being processed",
		},
		[]string{"endpoint"},
	)

	// 直方图：记录处理请求的时间
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)
