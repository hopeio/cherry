package gateway

import (
	"context"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"net/http"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hopeio/cherry/utils/encoding/protobuf/jsonpb"
	"google.golang.org/grpc/metadata"
)

type GatewayHandler func(context.Context, *runtime.ServeMux)

func Gateway(gatewayHandle GatewayHandler) *runtime.ServeMux {
	ctx := context.Background()

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &jsonpb.JSONPb{}),

		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			area, err := url.PathUnescape(req.Header.Get(httpi.HeaderArea))
			if err != nil {
				area = ""
			}
			var token = httpi.GetToken(req)

			return map[string][]string{
				httpi.HeaderArea:          {area},
				httpi.HeaderDeviceInfo:    {req.Header.Get(httpi.HeaderDeviceInfo)},
				httpi.HeaderLocation:      {req.Header.Get(httpi.HeaderLocation)},
				httpi.HeaderAuthorization: {token},
			}
		}),
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case
				"Accept",
				"Accept-Charset",
				"Accept-Language",
				"Accept-Ranges",
				//"Token",
				"Cache-Control",
				"Content-Type",
				//"Cookie",
				"Date",
				"Expect",
				"From",
				"Host",
				"If-Match",
				"If-Modified-Since",
				"If-None-Match",
				"If-Schedule-Tag-Match",
				"If-Unmodified-Since",
				"Max-Forwards",
				"Origin",
				"Pragma",
				"Referer",
				"User-Agent",
				"Via",
				"Warning":
				return key, true
			}
			return "", false
		}),
		runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcher))

	runtime.WithForwardResponseOption(CookieHook)(gwmux)
	runtime.WithForwardResponseOption(ResponseHook)(gwmux)
	runtime.WithRoutingErrorHandler(RoutingErrorHandler)(gwmux)
	runtime.WithErrorHandler(CustomHTTPError)(gwmux)
	if gatewayHandle != nil {
		gatewayHandle(ctx, gwmux)
	}
	return gwmux
}
