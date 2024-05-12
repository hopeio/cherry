package server

import (
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/utils/log"
	prometheus1 "github.com/hopeio/cherry/utils/net/http/prometheus"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

/*func init() {
	sink, _ := prometheus.NewPrometheusSink()
	conf := metrics.DefaultConfig(initialize.GlobalConfig.Module)
	metrics1, _ := metrics.New(conf, sink)
	metrics1.EnableHostnameLabel = true
	http.Handle("/metrics", promhttp.Handler())
}*/

var (
	meter      = otel.Meter("service-meter")
	apiCounter metric.Int64Counter
	histogram  metric.Float64Histogram
)

func init() {
	var err error
	apiCounter, err = meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		log.Fatal(err)
	}
	histogram, err = meter.Float64Histogram(
		"task.duration",
		metric.WithDescription("The duration of task execution."),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

type MetricsRecord = func(ctxi *httpctx.Context, uri, method string, code int)

var defaultMetricsRecord = func(ctxi *httpctx.Context, uri, method string, code int) {
	labels := prometheus2.Labels{
		"method": method,
		"uri":    uri,
	}
	t := time.Now().Sub(ctxi.RequestAt.Time)
	prometheus1.AccessCounter.With(labels).Add(1)
	prometheus1.QueueGauge.With(labels).Set(1)
	prometheus1.HttpDurationsHistogram.With(labels).Observe(float64(t) / 1000)
	prometheus1.HttpDurations.With(labels).Observe(float64(t) / 1000)
}

func SetMetricsRecord(metricsRecord MetricsRecord) {
	if metricsRecord != nil {
		defaultMetricsRecord = metricsRecord
	}
}
