package entrypoint

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kei6u/dogfood/driver"
	protov1 "github.com/kei6u/dogfood/proto/v1"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func RunBackend() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	startAPM()
	defer stopAPM()

	db, closeDB, err := driver.NewPsql()
	if err != nil {
		logger.Fatal("exit due to connection failure of database", zap.Error(err))
	}
	defer closeDB()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	go generateSpan(ctx, logger, "web")
	go generateSpan(ctx, logger, "http")
	go generateSpan(ctx, logger, "db")
	go generateSpan(ctx, logger, "cache")
	go generateSpan(ctx, logger, "function")

	s, err := protov1.NewServer(
		ctx,
		os.Getenv("GRPC_ADDR"),
		os.Getenv("GRPC_GATEWAY_ADDR"),
		logger,
		db,
	)
	if err != nil {
		logger.Fatal("exit due to a failure of initializeing dogfood backend server", zap.Error(err))
	}
	if err := s.Start(ctx); err != nil {
		logger.Warn("exit due to a failure of starting dogfood backend server", zap.Error(err))
		return
	}

	<-ctx.Done()
	s.Stop()
}

func generateSpan(ctx context.Context, logger *zap.Logger, spanType string) {
	spanName := fmt.Sprintf("%s_type_span_generator", spanType)
	for {
		if ctx.Err() != nil {
			return
		}
		span := tracer.StartSpan(spanName, tracer.SpanType(spanType))
		time.Sleep(time.Duration(rand.Intn(3) * int(time.Second)))
		logger.Info(fmt.Sprintf("%s span generated at %v", spanType, time.Now()), zap.Uint64("dd.span_id", span.Context().SpanID()), zap.Uint64("dd.trace_id", span.Context().TraceID()))
		span.Finish()
	}
}
