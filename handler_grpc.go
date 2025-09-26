/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"encoding/json"
	"fmt"

	"reflect"
	"runtime/debug"

	"github.com/hopeio/context/httpctx"
	grpcx "github.com/hopeio/gox/net/http/grpc"

	"github.com/hopeio/gox/log"
	runtimex "github.com/hopeio/gox/runtime"
	"github.com/hopeio/gox/validation/validator"
	"github.com/modern-go/reflect2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func (s *Server) grpcHandler() *grpc.Server {
	//conf := s.Config
	grpclog.SetLoggerV2(zapgrpc.NewLogger(log.CallerSkipLogger(4).Logger))
	if s.GrpcHandler != nil {
		var stream = append([]grpc.StreamServerInterceptor{StreamAccess, StreamValidator}, s.Grpc.StreamServerInterceptors...)
		var unary = append([]grpc.UnaryServerInterceptor{s.UnaryAccess, UnaryValidator}, s.Grpc.UnaryServerInterceptors...)
		// 想做的大而全几乎不可能,为了更高的自由度,这里不做实现,均由使用者自行实现,后续可提供默认实现,但同样要由用户自己调用
		/*		var srvMetrics *grpcprom.ServerMetrics
				if conf.EnableMetrics {
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

		s.Grpc.Options = append([]grpc.ServerOption{
			grpc.ChainStreamInterceptor(stream...),
			grpc.ChainUnaryInterceptor(unary...),
		}, s.Grpc.Options...)
		if s.Telemetry.Enabled {
			s.Grpc.Options = append(s.Grpc.Options, grpc.StatsHandler(otelgrpc.NewServerHandler(s.Telemetry.otelgrpcOpts...)))
		}

		grpcServer := grpc.NewServer(s.Grpc.Options...)
		/*		if conf.EnabledMetrics {
				srvMetrics.InitializeMetrics(grpcServer)
			}*/
		s.GrpcHandler(grpcServer)
		reflection.Register(grpcServer)
		return grpcServer
	}
	return nil
}

func (s *Server) UnaryAccess(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//enabledPrometheus := conf.EnabledMetrics

	defer func() {
		if r := recover(); r != nil {
			frame := debug.Stack()
			log.Errorw(fmt.Sprintf("panic: %v", r), zap.ByteString(log.FieldStack, frame))
			err = grpcx.Internal.ErrRep()
		}
	}()

	resp, err = handler(ctx, req)
	var code int
	//不能添加错误处理，除非所有返回的结构相同
	if err != nil {
		if v, ok := err.(interface{ GRPCStatus() *status.Status }); !ok {
			err = grpcx.Unknown.Msg(err.Error())
			code = int(grpcx.Unknown)
		} else {
			code = int(v.GRPCStatus().Code())
		}
	}
	if err == nil && reflect2.IsNil(resp) {
		resp = reflect.New(reflect.TypeOf(resp).Elem()).Interface()
	}
	/*		body, _ := protojson.Marshal(req.(proto.Msg)) // 性能比标准库差很多
			result, _ := protojson.Marshal(resp.(proto.Msg))*/
	body, _ := json.Marshal(req)
	result, _ := json.Marshal(resp)
	ctxi, _ := httpctx.FromContext(ctx)
	if s.HttpOption.AccessLog != nil {
		s.HttpOption.AccessLog(ctxi, &AccessLogParam{
			Method: "grpc",
			Url:    info.FullMethod,
			ReqBody: Body{
				IsJson: true,
				Data:   body,
			},
			RespBody: Body{
				IsJson: true,
				Data:   result,
			},
			StatusCode: code,
		})
	}
	/*		if enabledPrometheus {
			defaultMetricsRecord(ctxi, info.FullMethod, "grpc", code)
		}*/
	return resp, err

}

func StreamAccess(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			frame, _ := runtimex.GetCallerFrame(2)
			log.Errorw(fmt.Sprintf("panic: %v", r), zap.String(log.FieldStack, fmt.Sprintf("%s:%d (%#x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function)))
			err = grpcx.Internal.ErrRep()
		}
		//不能添加错误处理，除非所有返回的结构相同
		if err != nil {
			if _, ok := err.(interface{ GRPCStatus() *status.Status }); !ok {
				err = grpcx.Unknown.Msg(err.Error())
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
		return grpcx.InvalidArgument.Msg(validator.TransError(err))
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
		return nil, grpcx.InvalidArgument.Msg(validator.TransError(err))
	}
	return handler(ctx, req)
}

func StreamValidator(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	wrapper := &recvWrapper{
		ServerStream: stream,
	}
	return handler(srv, wrapper)
}
