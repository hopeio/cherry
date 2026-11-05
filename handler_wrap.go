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

// UnWrapContext returns what WrapContext put on ctx, or ctx itself when nothing was
// wrapped there.
//
// Falling back to ctx instead of nil lets a service hold one variable across the faces
// it is reached through: the gin wrapper hands back the *gin.Context (which implements
// context.Context and delegates to the request), gRPC and a plain http handler hand back
// the original ctx — either one still works as a context. Callers needing the concrete
// payload still type-assert, and a missing ReqResp now fails that assertion instead of
// returning a nil that panics on the next call.
func UnWrapContext(ctx context.Context) any {
	if v := ctx.Value(wrapContextKey); v != nil {
		return v
	}
	return ctx
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
