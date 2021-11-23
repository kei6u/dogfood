package grpc_backend

import (
	"context"
	"strings"

	"github.com/kei6u/dogfood/pkg/ddconfig"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func ddtracerUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if strings.Contains(info.FullMethod, "healthcheck") {
			return handler(ctx, req)
		}
		span, c := tracer.StartSpanFromContext(
			ctx,
			ddconfig.GetService(ddconfig.WithServiceSuffix(".grpcserver")),
		)
		defer span.Finish()
		return handler(c, req)
	}
}
