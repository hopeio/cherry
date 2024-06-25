package debug

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

func StackHandler() http.Handler {
	return http.HandlerFunc(Stack)
}

func Stack(w http.ResponseWriter, r *http.Request) {
	w.Write(debug.Stack())
}

func init() {
	http.Handle("/debug/", http.HandlerFunc(Stack))
}
