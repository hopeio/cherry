package cherry

// Deprecated 使用opentelemetry
import (
	"github.com/hopeio/context/httpctx"
	prometheus1 "github.com/hopeio/utils/net/http/prometheus"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"time"
)

/*func init() {
	sink, _ := prometheus.NewPrometheusSink()
	conf := metrics.DefaultConfig(initialize.GlobalConfig.Module)
	metrics1, _ := metrics.New(conf, sink)
	metrics1.EnableHostnameLabel = true
	http.Handle("/metrics", promhttp.Handler())
}*/

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
