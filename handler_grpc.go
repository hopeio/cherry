/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package mix

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hopeio/gox/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/stats/opentelemetry"
	"google.golang.org/grpc/status"
)

const (
	HeaderGrpcTraceBin = "Grpc-Trace-Bin"
	HeaderGrpcStatus   = "Grpc-Status"
	HeaderGrpcMsg      = "Grpc-Msg"
	HeaderGrpcMessage  = "Grpc-Message"
)


type GRPCStatus interface {
	GRPCStatus() *status.Status
}

var grpcLoggerOnce sync.Once

func (s *Server) grpcHandler() *grpc.Server {
	grpcLoggerOnce.Do(func() {
		grpclog.SetLoggerV2(zapgrpc.NewLogger(log.NoCallerLogger().With(zap.String("server", "grpc")).Logger))
	})
	if s.GrpcHandler != nil {
		var stream []grpc.StreamServerInterceptor
		var unary []grpc.UnaryServerInterceptor

		if s.Otel.Enabled {
			s.Grpc.Options = append(s.Grpc.Options, grpc.StatsHandler(
				otelgrpc.NewServerHandler(append([]otelgrpc.Option{
					otelgrpc.WithPropagators(propagation.NewCompositeTextMapPropagator(
						propagation.TraceContext{},
						opentelemetry.GRPCTraceBinPropagator{},
						propagation.Baggage{},
					)),
					// Treat calls that fail the internal auth check as public
					// (untrusted) endpoints: start a fresh root span instead of
					// inheriting the client trace, and link the client's span
					// context for correlation. Internal service-to-service calls
					// carry Server.InternalAuth* (shared secret) and are trusted —
					// their traceparent is inherited normally.
					otelgrpc.WithPublicEndpointFn(func(ctx context.Context, _ *stats.RPCTagInfo) bool {
						return !s.IsInternalCall(ctx)
					}),
				}, s.Otel.OtelgrpcOpts...)...)))
		}
		stream = append(stream, s.StreamAccess)
		stream = append(stream, s.Grpc.StreamServerInterceptors...)
		unary = append(unary, s.UnaryAccess)
		unary = append(unary, s.Grpc.UnaryServerInterceptors...)

		s.Grpc.Options = append(s.Grpc.Options, grpc.ChainStreamInterceptor(stream...),
			grpc.ChainUnaryInterceptor(unary...))

		grpcServer := grpc.NewServer(s.Grpc.Options...)
		s.GrpcHandler(grpcServer)
		reflection.Register(grpcServer)
		// The standard health service, for grpc-health-probe and the k8s gRPC probe.
		healthgrpc.RegisterHealthServer(grpcServer, health.NewServer())
		return grpcServer
	}
	return nil
}

func (s *Server) UnaryAccess(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//enabledPrometheus := conf.EnabledMetrics

	defer func() {
		if r := recover(); r != nil {
			log.StackLogger().Errorw(fmt.Sprintf("panic: %v", r))
			err = status.Error(codes.Internal, sysErrMsg)
		}
	}()
	// A dedicated gRPC port does not go through the HTTP dispatch layer, so
	// the context may carry no Metadata; inject a fallback here.
	md := GetMetadata(ctx)
	if md == nil {
		md = &Metadata{RequestAt: time.Now()}
		ctx = WithMetadata(ctx, md)
	}
	md.TraceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	md.Logger = log.DefaultLogger().With(zap.String(log.FieldTraceId, md.TraceId))
	md.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	var ok bool
	md.IncomingMD, ok = metadata.FromIncomingContext(ctx)
	if !ok {
		md.IncomingMD = nil
	}
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		md.PeerAddr = p.Addr.String()
	}
	if err = ValidateStruct(req); err != nil {
		return nil, err
	}
	resp, err = handler(ctx, req)

	if err != nil {
		if _, ok := err.(GRPCStatus); !ok {
			log.Errorw("untyped grpc handler error",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			err = Internal.ErrResp()
		}
	}

	// Tolerate a handler that returns two nils: the gRPC proto codec cannot
	// marshal a nil response and reporting it as is fails with Internal, so
	// build an empty response from the method signature instead.
	if err == nil && resp == nil {
		resp = emptyResp(info)
	}

	if s.Grpc.RecordFunc != nil {
		s.Grpc.RecordFunc(ctx, &GrpcAccessLogParam{
			Method:   info.FullMethod,
			Metadata: md,
			Request:  req,
			Response: resp,
			Err:      err,
		})
	}
	return resp, err
}

// emptyRespTypeCache caches FullMethod -> response element type, so that
// every request does not have to reflect over the method signature.
var emptyRespTypeCache sync.Map

// emptyResp reflects the response type out of the service method signature
// and builds an empty instance of it.
// It cannot reflect on the nil resp itself (reflect.TypeOf(nil) == nil and
// calling Elem on it panics), so the type has to come from the method
// signature of info.Server. The type is cached per FullMethod: once it is
// warm, a request costs one sync.Map lookup and one small allocation.
func emptyResp(info *grpc.UnaryServerInfo) any {
	if v, ok := emptyRespTypeCache.Load(info.FullMethod); ok {
		if t, _ := v.(reflect.Type); t != nil {
			return reflect.New(t).Interface()
		}
		return nil
	}
	var t reflect.Type
	if i := strings.LastIndexByte(info.FullMethod, '/'); i >= 0 && i < len(info.FullMethod)-1 {
		if m, ok := reflect.TypeOf(info.Server).MethodByName(info.FullMethod[i+1:]); ok {
			ft := m.Type
			// The generated service method has the signature
			// (ctx, req) (*Resp, error); take the first return value that is
			// not an error.
			if ft.NumOut() > 0 && ft.Out(0) != reflect.TypeFor[error]() {
				rt := ft.Out(0)
				if rt.Kind() == reflect.Pointer {
					t = rt.Elem()
				} else {
					t = rt
				}
			}
		}
	}
	emptyRespTypeCache.Store(info.FullMethod, t)
	if t == nil {
		return nil
	}
	return reflect.New(t).Interface()
}

func (s *Server) StreamAccess(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.StackLogger().Errorw(fmt.Sprintf("panic: %v", r))
			err = status.Error(codes.Internal, sysErrMsg)
		}
	}()
	ctx := stream.Context()
	// A dedicated gRPC port does not go through the HTTP dispatch layer, so
	// the context may carry no Metadata; inject a fallback here.
	md := GetMetadata(ctx)
	if md == nil {
		md = &Metadata{RequestAt: time.Now()}
		ctx = WithMetadata(ctx, md)
	}
	md.TraceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	md.Logger = log.DefaultLogger().With(zap.String(log.FieldTraceId, md.TraceId))
	md.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	var ok bool
	md.IncomingMD, ok = metadata.FromIncomingContext(ctx)
	if !ok {
		md.IncomingMD = nil
	}
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		md.PeerAddr = p.Addr.String()
	}
	wrapper := &recvWrapper{
		ServerStream: stream,
		ctx:          ctx,
	}
	err = handler(srv, wrapper)
	if err != nil {
		if _, ok := err.(GRPCStatus); !ok {
			log.Errorw("untyped grpc stream handler error",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			err = Internal.ErrResp()
		}
	}

	if s.Grpc.RecordFunc != nil {
		wrapper.Err = err
		wrapper.Method = info.FullMethod
		s.Grpc.RecordFunc(wrapper.Context(), &wrapper.GrpcAccessLogParam)
	}
	return err
}

type recvWrapper struct {
	grpc.ServerStream
	GrpcAccessLogParam
	ctx context.Context
}

// Context returns the context with the Metadata injected, so that the
// business handler can reach it.
func (s *recvWrapper) Context() context.Context {
	return s.ctx
}

func (s *recvWrapper) SendMsg(m interface{}) error {
	s.Response = m
	return s.ServerStream.SendMsg(m)
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	// The data has to be received before it is validated, otherwise the
	// validation would run against a zero-valued struct.
	s.Request = m
	if err := ValidateStruct(m); err != nil {
		return err
	}
	return nil
}
