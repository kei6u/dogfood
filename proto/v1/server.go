package protov1

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kei6u/dogfood/pkg/ddconfig"
	dogfoodpb "github.com/kei6u/dogfood/proto/v1/dogfood"
	healthcheckpb "github.com/kei6u/dogfood/proto/v1/healthcheck"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpc_dd "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	http_dd "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var _ dogfoodpb.DogFoodServiceServer = (*Server)(nil)
var _ healthcheckpb.HealthCheckServiceServer = (*Server)(nil)

type Server struct {
	gRPCAddr       string
	gRPCGWAddr     string
	logger         *zap.Logger
	db             *sql.DB
	promMetrics    *grpc_prometheus.ServerMetrics
	promHttpServer *http.Server
	grpcServer     *grpc.Server
	grpcgwServer   *http.Server
	conngRPCServer *grpc.ClientConn
	gRPCListener   net.Listener
}

func NewServer(ctx context.Context, gRPCAddr, gRPCGWAddr string, logger *zap.Logger, db *sql.DB) (*Server, error) {
	if !strings.HasPrefix(gRPCAddr, ":") {
		gRPCAddr = fmt.Sprintf(":%s", gRPCAddr)
	}
	if !strings.HasPrefix(gRPCGWAddr, ":") {
		gRPCGWAddr = fmt.Sprintf(":%s", gRPCGWAddr)
	}
	s := &Server{
		gRPCAddr:    gRPCAddr,
		gRPCGWAddr:  gRPCGWAddr,
		logger:      logger,
		db:          db,
		promMetrics: grpc_prometheus.NewServerMetrics(),
	}
	if err := s.initializePromHttpServer(); err != nil {
		return nil, err
	}
	s.initializegRPCServer()
	if err := s.initializegRPCGatewayServer(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("dogfood backend server has started")

	lis, err := net.Listen("tcp", s.gRPCAddr)
	if err != nil {
		return fmt.Errorf("faild to listen to gRPC Address: %w", err)
	}
	s.gRPCListener = lis

	go s.promHttpServer.ListenAndServe()
	go s.grpcServer.Serve(lis)
	go s.grpcgwServer.ListenAndServe()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop() {
	ctx := context.Background()
	s.logger.Info("dogfood gateway is shutting down")
	if err := s.gRPCListener.Close(); err != nil {
		s.logger.Error("failed to close listener to gRPC server", zap.Error(err))
	}
	if err := s.promHttpServer.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown prometheus server", zap.Error(err))
	}
	if err := s.grpcgwServer.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown gRPC gateway server", zap.Error(err))
	}
	s.grpcServer.GracefulStop()
	s.logger.Info("bye~~")
}

func (s *Server) initializePromHttpServer() error {
	r := prometheus.NewRegistry()
	if err := r.Register(s.promMetrics); err != nil {
		return fmt.Errorf("failed to initialize gRPC metrics to prometheus: %w", err)
	}
	if err := r.Register(dogfoodGramGuage); err != nil {
		return fmt.Errorf("failed to initialize dogfood gram guage to prometheus: %w", err)
	}
	if err := r.Register(dogfoodNameCount); err != nil {
		return fmt.Errorf("failed to initialize dogfood name count to prometheus: %w", err)
	}
	s.promHttpServer = &http.Server{
		Handler: promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
		Addr:    ":9092",
	}
	return nil
}

func (s *Server) initializegRPCServer() {
	ignoreMethods := make([]string, len(healthcheckpb.HealthCheckService_ServiceDesc.Methods))
	for i, m := range healthcheckpb.HealthCheckService_ServiceDesc.Methods {
		ignoreMethods[i] = fmt.Sprintf(
			"/%s/%s",
			healthcheckpb.HealthCheckService_ServiceDesc.ServiceName,
			m.MethodName,
		)
	}

	grpcsvc := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_dd.UnaryServerInterceptor(
				grpc_dd.WithIgnoredMethods(ignoreMethods...),
			),
			metricsUnaryServerInterceptor(),
			s.promMetrics.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(
				s.logger,
				grpc_zap.WithDecider(func(fullMethodName string, _ error) bool {
					return !strings.Contains(fullMethodName, "healthcheck")
				}),
				grpc_zap.WithMessageProducer(func(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {
					if dds, ok := tracer.SpanFromContext(ctx); ok {
						grpc_zap.AddFields(
							ctx,
							zap.Uint64("dd.trace_id", dds.Context().TraceID()),
							zap.Uint64("dd.span_id", dds.Context().SpanID()),
						)
					}
					grpc_zap.DefaultMessageProducer(ctx, msg, level, code, err, duration)
				}),
			),
		),
	)
	dogfoodpb.RegisterDogFoodServiceServer(grpcsvc, s)
	healthcheckpb.RegisterHealthCheckServiceServer(grpcsvc, s)
	s.promMetrics.InitializeMetrics(grpcsvc)

	s.grpcServer = grpcsvc
}

func (s *Server) initializegRPCGatewayServer(ctx context.Context) error {
	conn, err := grpc.DialContext(
		ctx,
		s.gRPCAddr,
		grpc.WithInsecure(),
		grpc.WithDisableHealthCheck(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC Server: %w", err)
	}
	s.conngRPCServer = conn

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				tracer.DefaultTraceIDHeader:  r.Header.Get(tracer.DefaultTraceIDHeader),
				tracer.DefaultParentIDHeader: r.Header.Get(tracer.DefaultParentIDHeader),
				tracer.DefaultPriorityHeader: r.Header.Get(tracer.DefaultPriorityHeader),
			})
		}),
	)
	if err := dogfoodpb.RegisterDogFoodServiceHandler(ctx, gwmux, conn); err != nil {
		return fmt.Errorf("failed to regiser handler: %w", err)
	}
	if err := healthcheckpb.RegisterHealthCheckServiceHandler(ctx, gwmux, conn); err != nil {
		return fmt.Errorf("failed to regiser handler: %w", err)
	}
	s.grpcgwServer = &http.Server{
		Addr: s.gRPCGWAddr,
		Handler: http_dd.WrapHandler(
			gwmux,
			ddconfig.GetService(ddconfig.WithServiceSuffix(".grpcgateway")),
			"",
			http_dd.WithIgnoreRequest(func(r *http.Request) bool {
				return strings.Contains(strings.ToLower(r.RequestURI), "healthcheck")
			}),
		),
	}
	return nil
}
