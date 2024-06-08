package debug

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func DebugHandler() http.Handler {
	http.Handle("/debug/", http.HandlerFunc(Debug))
	return http.DefaultServeMux
}

func Debug(w http.ResponseWriter, r *http.Request) {
	w.Write(debug.Stack())
}

func PromHandler() http.Handler {
	http.Handle("/metrics", promhttp.Handler())
	return http.DefaultServeMux
}
