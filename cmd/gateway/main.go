package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	redisrate "github.com/go-redis/redis_rate/v9"
	"github.com/kei6u/dogfood/driver"
	"github.com/kei6u/dogfood/pkg/ddconfig"
	"go.uber.org/zap"
	http_dd "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

const (
	createRecordRequestURI   = "/v1/dogfood/record"
	listRecordsRequestURI    = "/v1/dogfood/records"
	livenessProbeRequestURI  = "/v1/healthcheck/livenessProbe"
	readinessProbeRequestURI = "/v1/healthcheck/readinessProbe"
	startupProbeRequestURI   = "/v1/healthcheck/startupProbe"
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
			createRecordRequestURI,
			listRecordsRequestURI,
		},
	); err != nil {
		logger.Fatal("failed to register a revere proxy to gateway", zap.Error(err))
	}

	// Reverse Proxy
	http.HandleFunc(gw.handleFunc(createRecordRequestURI))
	http.HandleFunc(gw.handleFunc(listRecordsRequestURI))

	// Health check
	http.HandleFunc(livenessProbeRequestURI, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	http.HandleFunc(readinessProbeRequestURI, func(w http.ResponseWriter, _ *http.Request) {
		if err := r.Ping(ctx).Err(); err != nil {
			logger.Error("readiness probe failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("readiness probe failed: %s", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc(startupProbeRequestURI, func(w http.ResponseWriter, _ *http.Request) {
		if err := r.Set(ctx, startupProbeRequestURI, true, time.Second).Err(); err != nil {
			logger.Error("startup probe failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("startup probe failed: %s", err)))
			return
		}
		if _, err := r.Get(ctx, startupProbeRequestURI).Result(); err != nil {
			logger.Error("startup probe failed", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("startup probe failed: %s", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	logger.Info("dogfood gateway has started", zap.String("port", addr))
	if err := (&http.Server{
		Addr: fmt.Sprintf(":%s", addr),
		Handler: http_dd.WrapHandler(
			http.DefaultServeMux,
			ddconfig.GetService(),
			"",
			http_dd.WithAnalytics(true),
			http_dd.WithIgnoreRequest(func(r *http.Request) bool {
				for _, uri := range []string{livenessProbeRequestURI, readinessProbeRequestURI, startupProbeRequestURI} {
					if strings.EqualFold(uri, r.RequestURI) {
						return true
					}
				}
				return false
			}),
		),
	}).ListenAndServe(); err != nil {
		logger.Warn("dogfood gateway fails to start", zap.Error(err))
		return
	}
}
