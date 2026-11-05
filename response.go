/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package mix

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"

	iox "github.com/hopeio/gox/io"
	httpx "github.com/hopeio/gox/net/http"
	strings "github.com/hopeio/gox/strings"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

type ResponseWriter interface {
	WriteHeader(code int)
	HeaderX() httpx.Header
	Write([]byte) (int, error)
}

type Responder interface {
	Respond(ctx context.Context, w http.ResponseWriter) (int, error)
}

var errCodeHttpStatusMap = map[ErrCode]int{
	Success: http.StatusOK,
}

func RegisterErrCodeHttpStatus(code ErrCode, status int) {
	errCodeHttpStatusMap[code] = status
}

type errHeaderKey struct{}

var ErrHeaderKey errHeaderKey

func RespodWithErrHeader(ctx context.Context) context.Context {
	return context.WithValue(ctx, ErrHeaderKey, true)
}

type headerSetter interface {
	Set(key, value string)
}

// WriteErrHeaders writes only the Error-Code header on a business error
// (the composite value: business code*100 + gRPC code), which is how the
// HTTP client tells that the call failed. The msg and the i18n variables
// always travel in the response body instead (the ErrorInfo of protobuf or
// the msg/data of JSON).
func WriteErrHeaders(h headerSetter, code ErrCode) {
	if h == nil || code == Success {
		return
	}
	h.Set(HeaderErrorCode, strconv.Itoa(int(code)))
}

func StatusFromErrCode(code ErrCode) int {
	if status, ok := errCodeHttpStatusMap[code]; ok {
		return status
	}
	return http.StatusOK
}

// CommonResp is meant for receiving a response; use CommonAnyResp to send
// one.
type CommonResp[T any] struct {
	Code ErrCode `json:"code"`
	Msg  string  `json:"msg,omitempty"`
	// Data carries the payload of the response.
	Data T `json:"data,omitempty"`
}

func (res *CommonResp[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.Respond(r.Context(), w)
}

func (res *CommonResp[T]) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	data, contentType, err := DefaultMarshal(ctx, res)
	if err != nil {
		return RespondError(ctx, w, err)
	}
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		WriteErrHeaders(header, res.Code)
		header.Set(httpx.HeaderContentType, contentType)
	} else {
		header := w.Header()
		WriteErrHeaders(header, res.Code)
		header.Set(httpx.HeaderContentType, contentType)
	}

	if res.Code != Success {
		w.WriteHeader(StatusFromErrCode(res.Code))
	} else {
		w.WriteHeader(http.StatusOK)
	}

	recordBody(w, data, res)
	return w.Write(data)
}

type CommonProtoResp[T proto.Message] CommonResp[T]

func NewCommonProtoResp[T proto.Message](code ErrCode, msg string, data T) *CommonProtoResp[T] {
	return &CommonProtoResp[T]{Code: code, Msg: msg, Data: data}
}

func (r *CommonProtoResp[T]) MarshalProto() ([]byte, error) {
	buf := make([]byte, 0, 64)

	if r.Code != 0 {
		buf = protowire.AppendVarint(buf, 0x08)
		buf = protowire.AppendVarint(buf, uint64(r.Code))
	}

	// Encode Msg (field 2, string).
	if r.Msg != "" {
		buf = protowire.AppendVarint(buf, 0x12)
		buf = protowire.AppendString(buf, r.Msg)
	}

	// Encode Data (field 3, length-delimited bytes) with standard protowire
	// tag 0x1a, whether or not Code is zero. The old encoder only wrote Data
	// when Code==0 and used a bare sentinel byte 0, so UnmarshalProto silently
	// dropped Data whenever Code!=0 or Msg!="".
	if !reflect.ValueOf(r.Data).IsNil() {
		data, err := proto.Marshal(r.Data)
		if err != nil {
			return nil, err
		}
		buf = protowire.AppendVarint(buf, 0x1a)
		buf = protowire.AppendVarint(buf, uint64(len(data)))
		buf = append(buf, data...)
	}

	return buf, nil
}

// UnmarshalProto decodes protobuf bytes into CommonProtoResp.
// Wire layout: field 1=Code (varint), 2=Msg (bytes), 3=Data (bytes).
func (r *CommonProtoResp[T]) UnmarshalProto(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var pos int
	for pos < len(data) {
		// Parse tag (field number + wire type).
		tag, n := protowire.ConsumeVarint(data[pos:])
		if n < 0 {
			return errors.New("invalid protobuf data: unable to consume varint")
		}
		pos += n

		fieldNum, wireType := protowire.DecodeTag(tag)
		switch fieldNum {
		case 1: // Code
			if wireType != protowire.VarintType {
				return errors.New("invalid wire type for Code field")
			}
			code, n := protowire.ConsumeVarint(data[pos:])
			if n < 0 {
				return errors.New("invalid protobuf data: unable to consume Code varint")
			}
			r.Code = ErrCode(code)
			pos += n

		case 2: // Msg
			if wireType != protowire.BytesType {
				return errors.New("invalid wire type for Msg field")
			}
			msg, n := protowire.ConsumeString(data[pos:])
			if n < 0 {
				return errors.New("invalid protobuf data: unable to consume Msg string")
			}
			r.Msg = msg
			pos += n

		case 3: // Data (length-delimited proto message)
			if wireType != protowire.BytesType {
				return errors.New("invalid wire type for Data field")
			}
			msg, n := protowire.ConsumeBytes(data[pos:])
			if n < 0 {
				return errors.New("invalid protobuf data: unable to consume Data bytes")
			}
			if reflect.ValueOf(r.Data).IsNil() {
				r.Data = r.Data.ProtoReflect().New().Interface().(T)
			}
			if err := proto.Unmarshal(msg, r.Data); err != nil {
				return err
			}
			pos += n

		default: // Skip unknown fields so unknown tags cannot loop forever.
			n := protowire.ConsumeFieldValue(fieldNum, wireType, data[pos:])
			if n < 0 {
				return errors.New("invalid protobuf data: unable to skip unknown field")
			}
			pos += n
		}
	}

	return nil
}

type CommonAnyResp = CommonResp[any]

func NewCommonAnyResp(code ErrCode, msg string, data any) *CommonAnyResp {
	return &CommonAnyResp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func ServeError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = Internal
	}
	ErrRespFrom(err).ServeHTTP(w, r)
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error) (int, error) {
	if err == nil {
		err = Internal
	}
	return ErrRespFrom(err).Respond(ctx, w)
}

func ServeSuccess(w http.ResponseWriter, r *http.Request, res any) {
	RespondSuccess(r.Context(), w, res)
}

func RespondSuccess(ctx context.Context, w http.ResponseWriter, res any) (int, error) {
	data, contentType, err := DefaultMarshal(ctx, res)
	if err != nil {
		return RespondError(ctx, w, err)
	}
	if wx, ok := w.(ResponseWriter); ok {
		wx.HeaderX().Set(httpx.HeaderContentType, contentType)
	} else {
		w.Header().Set(httpx.HeaderContentType, contentType)
	}
	recordBody(w, data, res)
	return w.Write(data)
}

// recordBody walks the Unwrap chain to find the RecordBodyer and records
// the response body: w may be the recorder itself, or it may be wrapped in
// several layers of middleware such as otelhttp.
func recordBody(w http.ResponseWriter, raw []byte, v any) {
	for ow := w; ; {
		if recorder, ok := ow.(httpx.RecordBodyer); ok {
			recorder.RecordBody(raw, v)
			return
		}
		uw, ok := ow.(httpx.Unwrapper)
		if !ok {
			return
		}
		ow = uw.Unwrap()
	}
}

func Serve(w http.ResponseWriter, r *http.Request, data any) {
	if err, ok := data.(error); ok {
		ServeError(w, r, err)
		return
	}
	ServeSuccess(w, r, data)
}

func Respond(ctx context.Context, w http.ResponseWriter, data any) (int, error) {
	if err, ok := data.(error); ok {
		return RespondError(ctx, w, err)
	}
	return RespondSuccess(ctx, w, data)
}

type Response struct {
	Status  int                `json:"status,omitempty"`
	Headers http.Header        `json:"header,omitempty"`
	Body    iox.WriterToCloser `json:"body,omitempty"`
}

func (res *Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.Respond(r.Context(), w)
}

func (res *Response) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		for k, v := range res.Headers {
			for _, vv := range v {
				header.Add(k, vv)
			}
		}
	} else {
		httpx.CopyHttpHeader(w.Header(), res.Headers)
	}
	statusCode := res.Status
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	w.WriteHeader(statusCode)
	n, err := res.Body.WriteTo(w)
	res.Body.Close()
	return int(n), err
}

type ErrResp CommonResp[map[string]string]

func NewErrResp(code ErrCode, msg string, data map[string]string) *ErrResp {
	return &ErrResp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func ErrRespFrom(err error) *ErrResp {
	if err == nil {
		return nil
	}
	if errresp, ok := err.(*ErrResp); ok {
		return errresp
	}
	type errresp interface{ ErrResp() *ErrResp }
	if se, ok := err.(errresp); ok {
		return se.ErrResp()
	}
	if st, ok := status.FromError(err); ok {
		er := &ErrResp{}
		// The msg is the composite value (business code*100 + gRPC code),
		// so it converts back directly.
		if c, perr := strconv.ParseInt(st.Message(), 10, 32); perr == nil {
			er.Code = ErrCode(c)
		} else {
			er.Code = ErrCode(st.Code())
		}
		for _, d := range st.Details() {
			if ei, ok := d.(*errdetails.ErrorInfo); ok {
				if ei.Reason != "" {
					er.Msg = ei.Reason
				}
				if len(ei.Metadata) > 0 {
					er.Data = ei.Metadata
				}
				break
			}
		}
		if er.Msg == "" {
			// Without an ErrorInfo, use the registered name; do not pass
			// status.Message down (it may be the driver's original text).
			er.Msg = er.Code.String()
		}
		return er
	}
	rv := reflect.ValueOf(err)
	kind := rv.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		c := ErrCode(rv.Int())
		return NewErrResp(c, c.String(), nil)
	}
	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		c := ErrCode(rv.Uint())
		return NewErrResp(c, c.String(), nil)
	}
	// An untyped error (the driver's original text): do not pass
	// err.Error() down.
	return NewErrResp(Internal, Internal.String(), nil)
}

func (res *ErrResp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.Respond(r.Context(), w)
}

func (res *ErrResp) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	if ctx != nil && ctx.Err() != nil {
		return 0, ctx.Err()
	}
	data, contentType, _ := DefaultMarshal(ctx, res)
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		WriteErrHeaders(header, res.Code)
		header.Set(httpx.HeaderContentType, contentType)
	} else {
		header := w.Header()
		WriteErrHeaders(header, res.Code)
		header.Set(httpx.HeaderContentType, contentType)
	}
	if res.Code != Success {
		w.WriteHeader(StatusFromErrCode(res.Code))
	} else {
		w.WriteHeader(http.StatusOK)
	}
	recordBody(w, data, res)
	return w.Write(data)
}

func (res *ErrResp) ErrResp() *ErrResp {
	return (*ErrResp)(res)
}

func (x *ErrResp) GRPCStatus() *status.Status {
	// The code uses the standard enum in the low two digits (below 100, so
	// it stays valid) and the msg carries the composite value, which the
	// client converts back into the business code. The i18n key and the {k}
	// placeholder variables go into the reason/metadata of the standard
	// google.rpc.ErrorInfo: it is the standard type for gRPC error details,
	// so the grpc-go and grpc-dart clients decode it on their own without
	// hand-written wire code.
	st := status.New(x.Code.GRPCCode(), strconv.Itoa(int(x.Code)))
	if x.Msg != "" || len(x.Data) > 0 {
		if st2, err := st.WithDetails(x.ErrorInfo()); err == nil {
			return st2
		}
	}
	return st
}

// ErrorInfo encodes the HTTP protobuf error body into a
// google.rpc.ErrorInfo: reason is the i18n key (the same as msg) and metadata
// holds the placeholder variables (the same as data).
func (x *ErrResp) ErrorInfo() *errdetails.ErrorInfo {
	if x == nil {
		return &errdetails.ErrorInfo{}
	}
	return &errdetails.ErrorInfo{
		Reason:   x.Msg,
		Metadata: x.Data,
	}
}

func (x *ErrResp) Error() string {
	return fmt.Sprintf("code: %d, msg: %s, data: %v", x.Code, x.Msg, x.Data)
}

type ResponseStream struct {
	Status  int                          `json:"status,omitempty"`
	Headers http.Header                  `json:"header,omitempty"`
	Body    iter.Seq[iox.WriterToCloser] `json:"body,omitempty"`
}

func (res *ResponseStream) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		for k, v := range res.Headers {
			for _, vv := range v {
				header.Add(k, vv)
			}
		}
	} else {
		httpx.CopyHttpHeader(w.Header(), res.Headers)
	}
	statusCode := res.Status
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	w.WriteHeader(statusCode)
	return RespondStream(ctx, w, res.Body)
}

func (res *ResponseStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.Respond(r.Context(), w)
}

func RespondStream(ctx context.Context, w http.ResponseWriter, dataSource iter.Seq[iox.WriterToCloser]) (int, error) {
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		header.Set(httpx.HeaderCacheControl, httpx.CacheControlNoCache)
		header.Set(httpx.HeaderTransferEncoding, httpx.TransferEncodingChunked)
	} else {
		header := w.Header()
		header.Set(httpx.HeaderCacheControl, httpx.CacheControlNoCache)
		header.Set(httpx.HeaderTransferEncoding, httpx.TransferEncodingChunked)
	}
	var n int
	flusher := w.(http.Flusher)
	for data := range dataSource {
		select {
		case <-ctx.Done():
			return n, ctx.Err()
		default:
			writen, err := data.WriteTo(w)
			n += int(writen)
			if err != nil {
				return n, err
			}
			flusher.Flush()
		}
	}
	return n, nil
}

type SSEData string

func (data SSEData) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(strings.ToBytes(fmt.Sprintf("data: %s\n\n", data)))
	return int64(n), err
}
func (data SSEData) Close() error {
	return nil
}

func RespondSSE[T ~string](ctx context.Context, w http.ResponseWriter, dataSource iter.Seq[T]) (int, error) {
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		header.Set(httpx.HeaderContentType, httpx.ContentTypeTextEventStream)
	} else {
		header := w.Header()
		header.Set(httpx.HeaderContentType, httpx.ContentTypeTextEventStream)
	}
	return RespondStream(ctx, w, func(yield func(iox.WriterToCloser) bool) {
		for data := range dataSource {
			yield(SSEData(data))
		}
	})
}

type XXXResponseBody interface {
	XXX_ResponseBody() any
}

type ResponseBody interface {
	ResponseBody() ([]byte, string)
}

type StatusCode interface {
	StatusCode(v any) int
}

type ResponseFile struct {
	Name        string             `json:"name"`
	Body        iox.WriterToCloser `json:"body"`
	ContentType string
}

func (res *ResponseFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.Respond(r.Context(), w)
}

func (res *ResponseFile) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	contentType := res.ContentType
	var limitedWriter iox.LimitedWriter
	if contentType == "" {
		if ext := filepath.Ext(res.Name); ext != "" {
			contentType = mime.TypeByExtension(ext)
			if contentType == "" {
				contentType = httpx.ContentTypeOctetStream
			}
		} else {
			limitedWriter = iox.NewLimitedWriter(512)
			// io.EOF means the probe buffer is full and the stream still has
			// more (the LimitedWriter convention); the rest goes out through the
			// second WriteTo below. The old behaviour silently dropped
			// everything beyond the 512 bytes.
			_, err := res.Body.WriteTo(&limitedWriter)
			if err != nil && !errors.Is(err, io.EOF) {
				return 0, err
			}
			contentType = http.DetectContentType(limitedWriter)
		}
	}
	contentDisposition := "inline"
	if res.Name != "" {
		contentDisposition = fmt.Sprintf(httpx.AttachmentTmpl, res.Name)
	}
	if wx, ok := w.(ResponseWriter); ok {
		header := wx.HeaderX()
		header.Set(httpx.HeaderContentType, contentType)
		header.Set(httpx.HeaderContentDisposition, contentDisposition)
	} else {
		header := w.Header()
		header.Set(httpx.HeaderContentType, contentType)
		header.Set(httpx.HeaderContentDisposition, contentDisposition)
	}
	if len(limitedWriter) > 0 {
		n, err := w.Write(limitedWriter)
		if err != nil {
			return n, err
		}
	}
	n, err := res.Body.WriteTo(w)
	res.Body.Close()
	if len(limitedWriter) > 0 {
		n += int64(len(limitedWriter))
	}
	return int(n), err
}
