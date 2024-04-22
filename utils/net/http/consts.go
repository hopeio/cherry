package http

const (
	// ContentBinaryHeaderValue header value for binary data.
	ContentBinaryHeaderValue = "application/octet-stream"
	// ContentWebassemblyHeaderValue header value for web assembly files.
	ContentWebassemblyHeaderValue = "application/wasm"
	// ContentHtmlHeaderValue is the  string of text/html response header's content type value.
	ContentHtmlHeaderValue = "text/html"
	// ContentJsonHeaderValue header value for JSON data.
	ContentJsonHeaderValue = "application/json"
	// ContentJsonProblemHeaderValue header value for JSON API problem error.
	// Read more at: https://tools.ietf.org/html/rfc7807
	ContentJsonProblemHeaderValue = "application/problem+json"
	// ContentXmlProblemHeaderValue header value for XML API problem error.
	// Read more at: https://tools.ietf.org/html/rfc7807
	ContentXmlProblemHeaderValue = "application/problem+xml"
	// ContentJavascriptHeaderValue header value for JSONP & Javascript data.
	ContentJavascriptHeaderValue = "text/javascript"
	// ContentTextHeaderValue header value for Text data.
	ContentTextHeaderValue = "text/plain"
	// ContentXmlHeaderValue header value for XML data.
	ContentXmlHeaderValue = "text/xml"
	// ContentXmlUnreadableHeaderValue obselete header value for XML.
	ContentXmlUnreadableHeaderValue = "application/xml"
	// ContentMarkdownHeaderValue custom key/content type, the real is the text/html.
	ContentMarkdownHeaderValue = "text/markdown"
	// ContentYamlHeaderValue header value for YAML data.
	ContentYamlHeaderValue = "application/x-yaml"
	// ContentYamlTextHeaderValue header value for YAML plain text.
	ContentYamlTextHeaderValue = "text/yaml"
	// ContentProtobufHeaderValue header value for Protobuf messages data.
	ContentProtobufHeaderValue = "application/x-protobuf"
	// ContentMsgPackHeaderValue header value for MsgPack data.
	ContentMsgPackHeaderValue = "application/msgpack"
	// ContentMsgPack2HeaderValue alternative header value for MsgPack data.
	ContentMsgPack2HeaderValue = "application/x-msgpack"
	// ContentFormHeaderValue header value for post form data.
	ContentFormHeaderValue = "application/x-www-form-urlencoded"
	// ContentFormMultipartHeaderValue header value for post multipart form data.
	ContentFormMultipartHeaderValue = "multipart/form-data"
	// ContentGrpcHeaderValue Content-Type header value for gRPC.
	ContentGrpcHeaderValue    = "application/grpc"
	ContentGrpcWebHeaderValue = "application/grpc-web"

	ContentJsonUtf8HeaderValue = "application/json;charset=utf-8"

	ContentFormParamHeaderValue = "application/x-www-form-urlencoded;param=value"
)

const (
	HeaderDeviceInfo                  = "Device-AuthInfo"
	HeaderLocation                    = "Location"
	HeaderArea                        = "Area"
	HeaderUserAgent                   = "User-Agent"
	HeaderXForwardedFor               = "X-Forwarded-For"
	HeaderAuth                        = "HeaderAuth"
	HeaderContentType                 = "Content-Type"
	HeaderTrace                       = "Tracing"
	HeaderTraceID                     = "Tracing-ID"
	HeaderTraceBin                    = "Tracing-Bin"
	HeaderAuthorization               = "Authorization"
	HeaderCookie                      = "Cookie"
	HeaderCookieToken                 = "token"
	HeaderCookieDel                   = "del"
	HeaderContentDisposition          = "Content-Disposition"
	HeaderContentEncoding             = "Content-Encoding"
	HeaderReferer                     = "Referer"
	HeaderAccept                      = "Accept"
	HeaderAcceptLanguage              = "Accept-Language"
	HeaderAcceptEncoding              = "Accept-Encoding"
	HeaderCacheControl                = "Cache-Control"
	HeaderSetCookie                   = "Set-Cookie"
	HeaderTrailer                     = "Trailer"
	HeaderTransferEncoding            = "Transfer-Encoding"
	HeaderInternal                    = "Internal"
	HeaderTE                          = "TE"
	HeaderLastModified                = "Last-Modified"
	HeaderContentLength               = "Content-Length"
	HeaderAccessControlRequestMethod  = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
	HeaderOrigin                      = "Origin"
	HeaderConnection                  = "Connection"
	HeaderRange                       = "Range"
	HeaderContentRange                = "Content-Range"
	HeaderAcceptRanges                = "Accept-Ranges"
)

const (
	HeaderGrpcTraceBin = "grpc-trace-bin"
	HeaderGrpcInternal = "grpc-internal"
)
