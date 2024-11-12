package cherry

// Deprecated 使用opentelemetry
import (
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/hopeio/context/httpctx"
	prometheus1 "github.com/hopeio/utils/net/http/prometheus"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"time"
)

/*func init() {
	sink, _ := prometheus.NewPrometheusSink()
	conf := metrics.DefaultConfig("")
	metrics1, _ := metrics.New(conf, sink)
	metrics1.EnableHostnameLabel = true
	http.Handle("/metrics", promhttp.Handler())
	reg.MustRegister(srvMetrics)
}*/

var reg = prometheus2.NewRegistry()

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

var srvMetrics = grpcprom.NewServerMetrics(
	grpcprom.WithServerHandlingTimeHistogram(
		grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
	),
)
