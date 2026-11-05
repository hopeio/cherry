package mix

import (
	"context"
	"net/http"

	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/types"
	"go.uber.org/zap"
)

type Service[REQ, RESP any] func(ctx ReqResp, req REQ) (RESP, *ErrResp)

type wrapKey struct{}

var wrapContextKey = wrapKey{}

// WrapContext puts the request/response pair on ctx for the service to pick
// up with UnWrapContext.
//
// It must derive from the request context. Starting a fresh
// context.Background() here would drop the trace span, the values added by
// middleware and the cancellation that follows a client disconnect or a
// server shutdown - the service would keep running against a detached root.
func WrapContext(ctx context.Context, v any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, wrapContextKey, v)
}

func UnWrapContext(ctx context.Context) any {
	return ctx.Value(wrapContextKey)
}

type ReqResp struct {
	*http.Request
	http.ResponseWriter
}

func HandlerWrap[REQ, RESP any](service Service[*REQ, *RESP]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := new(REQ)
		err := Bind(r, req)
		if err != nil {
			ServeError(w, r, InvalidArgument.Wrap(err))
			return
		}
		res, resperr := service(ReqResp{r, w}, req)
		if resperr != nil {
			ServeError(w, r, resperr)
			return
		}
		switch httpres := any(res).(type) {
		case http.Handler:
			httpres.ServeHTTP(w, r)
			return
		case Responder:
			httpres.Respond(r.Context(), w)
			return
		}
		ServeSuccess(w, r, res)
	})
}

func HandlerWrapCommon[REQ, RESP any](service types.Service[*REQ, *RESP]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := new(REQ)
		err := Bind(r, req)
		if err != nil {
			ServeError(w, r, InvalidArgument.Wrap(err))
			return
		}
		res, err := service(WrapContext(r.Context(), ReqResp{r, w}), req)
		if err != nil {
			if _, ok := err.(GRPCStatus); !ok {
				log.Errorw("untyped http handler error", zap.Error(err))
			}
			ErrRespFrom(err).ServeHTTP(w, r)
			return
		}
		switch httpres := any(res).(type) {
		case http.Handler:
			httpres.ServeHTTP(w, r)
			return
		case Responder:
			httpres.Respond(r.Context(), w)
			return
		}
		ServeSuccess(w, r, res)
	})
}
