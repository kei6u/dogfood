package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	redisrate "github.com/go-redis/redis_rate/v9"
	"github.com/kei6u/dogfood/driver"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

const (
	createRecordEndpoint   = "/v1/dogfood/record"
	listRecordsEndpoint    = "/v1/dogfood/records"
	livenessProbeEndpoint  = "/v1/healthcheck/livenessProbe"
	readinessProbeEndpoint = "/v1/healthcheck/readinessProbe"
	startupProbeEndpoint   = "/v1/healthcheck/startupProbe"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tracer.Start(tracer.WithAnalytics(true))
	defer tracer.Stop()

	if err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	); err != nil {
		logger.Warn("failed to start profiler", zap.Error(err))
		return
	}
	defer profiler.Stop()

	// Loading environment variables.
	var addr string
	var dogfoodBackendAddr string
	if addr = os.Getenv("ADDR"); addr == "" {
		logger.Fatal("ADDR is missing")
	}
	if dogfoodBackendAddr = os.Getenv("DOGFOOD_BACKEND_ADDR"); dogfoodBackendAddr == "" {
		logger.Fatal("DOGFOOD_BACKEND_ADDR is missing")
	}

	ctx := context.Background()

	// Creating connection with Redis.
	r, rClose, err := driver.NewRedis(ctx)
	if err != nil {
		logger.Fatal("failed to initialize redis client", zap.Error(err))
	}
	defer rClose()
	limiter := redisrate.NewLimiter(r)

	gw := newGateway(limiter, logger)
	if err := gw.registerReverseProxy(
		dogfoodBackendAddr,
		[]string{
			createRecordEndpoint,
			listRecordsEndpoint,
		},
	); err != nil {
		logger.Fatal("failed to register a revere proxy to gateway", zap.Error(err))
	}

	// Reverse Proxy
	http.HandleFunc(gw.handleFunc(createRecordEndpoint))
	http.HandleFunc(gw.handleFunc(listRecordsEndpoint))

	// Health check
	http.HandleFunc(livenessProbeEndpoint, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	http.HandleFunc(readinessProbeEndpoint, func(w http.ResponseWriter, _ *http.Request) {
		if err := r.Ping(ctx).Err(); err != nil {
			logger.Error("readiness probe failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("readiness probe failed: %s", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc(startupProbeEndpoint, func(w http.ResponseWriter, _ *http.Request) {
		if err := r.Set(ctx, startupProbeEndpoint, true, time.Second).Err(); err != nil {
			logger.Error("startup probe failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("startup probe failed: %s", err)))
			return
		}
		if _, err := r.Get(ctx, startupProbeEndpoint).Result(); err != nil {
			logger.Error("startup probe failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("startup probe failed: %s", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	logger.Info("gateway has started", zap.String("port", addr))
	http.ListenAndServe(fmt.Sprintf(":%s", addr), nil)
}
