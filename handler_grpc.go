/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"fmt"

	"reflect"
	"runtime/debug"

	"github.com/hopeio/gox/context/httpctx"
	grpcx "github.com/hopeio/gox/net/http/grpc"

	"github.com/hopeio/gox/log"
	runtimex "github.com/hopeio/gox/runtime"
	"github.com/hopeio/gox/validator"
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

type GRPCStatus interface {
	GRPCStatus() *status.Status
}

func (s *Server) UnaryAccess(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//enabledPrometheus := conf.EnabledMetrics

	defer func() {
		if r := recover(); r != nil {
			frame := debug.Stack()
			log.Errorw(fmt.Sprintf("panic: %v", r), zap.ByteString(log.FieldStack, frame))
			err = grpcx.Internal.ErrResp()
		}
	}()

	resp, err = handler(ctx, req)
	//不能添加错误处理，除非所有返回的结构相同
	if err != nil {
		if _, ok := err.(GRPCStatus); !ok {
			err = grpcx.Unknown.Msg(err.Error())
		}
	}
	if err == nil && reflect2.IsNil(resp) {
		resp = reflect.New(reflect.TypeOf(resp).Elem()).Interface()
	}

	ctxi, _ := httpctx.FromContext(ctx)
	if s.Grpc.RecordFunc != nil {
		s.Grpc.RecordFunc(ctxi, &GrpcAccessLogParam{
			Method: info.FullMethod,
			req:    req,
			resp:   resp,
			err:    err,
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
			err = grpcx.Internal.ErrResp()
		}
		//不能添加错误处理，除非所有返回的结构相同
		if err != nil {
			if _, ok := err.(GRPCStatus); !ok {
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
	if err := validator.ValidateStruct(m); err != nil {
		return grpcx.InvalidArgument.Wrap(err)
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

	if err = validator.ValidateStruct(req); err != nil {
		return nil, grpcx.InvalidArgument.Wrap(err)
	}
	return handler(ctx, req)
}

func StreamValidator(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	wrapper := &recvWrapper{
		ServerStream: stream,
	}
	return handler(srv, wrapper)
}
