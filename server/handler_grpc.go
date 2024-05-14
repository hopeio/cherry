package server

import (
	"context"
	"encoding/json"
	"fmt"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/log"
	runtimei "github.com/hopeio/cherry/utils/runtime"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/hopeio/cherry/utils/validation/validator"
	"github.com/modern-go/reflect2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"reflect"
	"runtime/debug"
)

func (s *Server) grpcHandler() *grpc.Server {
	//conf := s.Config
	grpclog.SetLoggerV2(zapgrpc.NewLogger(log.GetCallerSkipLogger(3).Logger))
	if s.GrpcHandler != nil {
		var stream = []grpc.StreamServerInterceptor{StreamAccess, StreamValidator}
		var unary = []grpc.UnaryServerInterceptor{UnaryAccess(s.Config), UnaryValidator}
		// 想做的大而全几乎不可能,为了更高的自由度,这里不做实现,均由使用者自行实现,后续可提供默认实现,但同样要由用户自己调用
		/*		var srvMetrics *grpcprom.ServerMetrics
				if conf.Metrics {
					// Setup metrics.
					srvMetrics = grpcprom.NewServerMetrics(
						grpcprom.WithServerHandlingTimeHistogram(
							grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
						),
					)
					prometheus.MustRegister(srvMetrics)
					exemplarFromContext := func(ctx context.Context) prometheus.Labels {
						if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
							return prometheus.Labels{"traceID": span.TraceID().String()}
						}
						return nil
					}
					stream = append(stream, srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)))
					unary = append(unary, srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)))
				}*/

		stream = append(stream, grpc_validator.StreamServerInterceptor())
		unary = append(unary, grpc_validator.UnaryServerInterceptor())
		s.Config.GrpcOptions = append([]grpc.ServerOption{
			grpc.ChainStreamInterceptor(stream...),
			grpc.ChainUnaryInterceptor(unary...),
		}, s.Config.GrpcOptions...)

		grpcServer := grpc.NewServer(s.Config.GrpcOptions...)
		/*		if conf.Metrics {
				srvMetrics.InitializeMetrics(grpcServer)
			}*/
		s.GrpcHandler(grpcServer)
		reflection.Register(grpcServer)
		return grpcServer
	}
	return nil
}

func UnaryAccess(conf *Config) grpc.UnaryServerInterceptor {
	enablePrometheus := conf.Metrics
	return func(
		ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				frame := debug.Stack()
				log.Errorw(fmt.Sprintf("panic: %v", r), zap.ByteString(log.FieldStack, frame))
				err = errorcode.SysError.ErrRep()
			}
		}()

		resp, err = handler(ctx, req)
		var code int
		//不能添加错误处理，除非所有返回的结构相同
		if err != nil {
			if v, ok := err.(interface{ GrpcStatus() *status.Status }); !ok {
				err = errorcode.Unknown.Message(err.Error())
				code = int(errorcode.Unknown)
			} else {
				code = int(v.GrpcStatus().Code())
			}
		}
		if err == nil && reflect2.IsNil(resp) {
			resp = reflect.New(reflect.TypeOf(resp).Elem()).Interface()
		}
		/*		body, _ := protojson.Marshal(req.(proto.Message)) // 性能比标准库差很多
				result, _ := protojson.Marshal(resp.(proto.Message))*/
		body, _ := json.Marshal(req)
		result, _ := json.Marshal(resp)
		ctxi := httpctx.ContextFromContext(ctx)
		defaultAccessLog(ctxi, info.FullMethod, "grpc",
			stringsi.BytesToString(body), stringsi.BytesToString(result),
			code)
		if enablePrometheus {
			defaultMetricsRecord(ctxi, info.FullMethod, "grpc", code)
		}
		return resp, err
	}
}

func StreamAccess(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			frame, _ := runtimei.GetCallerFrame(2)
			log.Errorw(fmt.Sprintf("panic: %v", r), zap.String(log.FieldStack, fmt.Sprintf("%s:%d (%#x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function)))
			err = errorcode.SysError.ErrRep()
		}
		//不能添加错误处理，除非所有返回的结构相同
		if err != nil {
			if _, ok := err.(interface{ GrpcStatus() *status.Status }); !ok {
				err = errorcode.Unknown.Message(err.Error())
			}
		}
	}()

	return handler(srv, stream)
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) SendMsg(m interface{}) error {
	return s.ServerStream.SendMsg(m)
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := validator.Validator.Struct(m); err != nil {
		return errorcode.InvalidArgument.Message(validator.Trans(err))
	}
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

func UnaryValidator(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	if err = validator.Validator.Struct(req); err != nil {
		return nil, errorcode.InvalidArgument.Message(validator.Trans(err))
	}
	return handler(ctx, req)
}

func StreamValidator(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	wrapper := &recvWrapper{
		ServerStream: stream,
	}
	return handler(srv, wrapper)
}
