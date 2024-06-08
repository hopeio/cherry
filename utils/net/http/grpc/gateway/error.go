package gateway

import (
	"fmt"
	"github.com/hopeio/cherry/protobuf/errorcode"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/cherry/utils/net/http/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/textproto"
)

func HTTPStatusFromCode(code errorcode.ErrCode) int {
	switch code {
	case errorcode.Success:
		return http.StatusOK
	case errorcode.Canceled:
		return http.StatusRequestTimeout
	case errorcode.Unknown:
		return http.StatusInternalServerError
	case errorcode.InvalidArgument:
		return http.StatusBadRequest
	case errorcode.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case errorcode.NotFound:
		return http.StatusNotFound
	case errorcode.AlreadyExists:
		return http.StatusConflict
	case errorcode.PermissionDenied:
		return http.StatusForbidden
	case errorcode.Unauthenticated:
		return http.StatusUnauthorized
	case errorcode.ResourceExhausted:
		return http.StatusTooManyRequests
	case errorcode.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case errorcode.Aborted:
		return http.StatusConflict
	case errorcode.OutOfRange:
		return http.StatusBadRequest
	case errorcode.Unimplemented:
		return http.StatusNotImplemented
	case errorcode.Internal:
		return http.StatusInternalServerError
	case errorcode.Unavailable:
		return http.StatusServiceUnavailable
	case errorcode.DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}

func OutgoingHeaderMatcher(key string) (string, bool) {
	switch key {
	case
		httpi.HeaderSetCookie:
		return key, true
	}
	return "", false
}

var headerMatcher = []string{httpi.HeaderSetCookie}

func HandleForwardResponseServerMetadata(w http.ResponseWriter, md metadata.MD) {
	for _, k := range headerMatcher {
		if vs, ok := md[k]; ok {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
	}
}

func HandleForwardResponseTrailerHeader(w http.ResponseWriter, md metadata.MD) {
	for k := range md {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", grpc.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func HandleForwardResponseTrailer(w http.ResponseWriter, md metadata.MD) {
	for k, vs := range md {
		tKey := fmt.Sprintf("%s%s", grpc.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}
