package entrypoint

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	redisrate "github.com/go-redis/redis_rate/v9"
	"github.com/kei6u/dogfood/driver"
	"github.com/kei6u/dogfood/pkg/ddconfig"
	"github.com/kei6u/dogfood/pkg/httplib"
	"go.uber.org/zap"
	http_dd "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	createRecordRequestURI   = "/v1/dogfood/record"
	listRecordsRequestURI    = "/v1/dogfood/records"
	livenessProbeRequestURI  = "/v1/healthcheck/livenessProbe"
	readinessProbeRequestURI = "/v1/healthcheck/readinessProbe"
	startupProbeRequestURI   = "/v1/healthcheck/startupProbe"
)

func RunGateway() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	startAPM()
	defer stopAPM()

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

	s := &http.Server{
		Addr: fmt.Sprintf(":%s", addr),
		Handler: http_dd.WrapHandler(
			http.DefaultServeMux,
			ddconfig.GetService(),
			"",
			http_dd.WithAnalytics(true),
			http_dd.WithIgnoreRequest(func(r *http.Request) bool {
				return strings.Contains(strings.ToLower(r.RequestURI), "healthcheck")
			}),
		),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		logger.Info("dogfood gateway has started", zap.String("port", addr))
		if err := s.ListenAndServe(); err != nil {
			logger.Info("failed to listen and serve", zap.Error(err))
		}
	}()
	go func() {
		<-c
		logger.Info("dogfood gateway is shutting down")
		s.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
}

type gateway struct {
	limiter    *redisrate.Limiter
	rpLookup   map[string]*httputil.ReverseProxy // key: addr, e.g. /v1/dogfood/record
	addrLookup map[string]string                 // key: addr, value: pattern
	l          *zap.Logger
}

func newGateway(limiter *redisrate.Limiter, l *zap.Logger) *gateway {
	return &gateway{
		limiter,
		make(map[string]*httputil.ReverseProxy),
		map[string]string{},
		l,
	}
}

func (gw *gateway) registerReverseProxy(addr string, patterns []string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return fmt.Errorf("target host is invalid: %w", err)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	for _, p := range patterns {
		gw.addrLookup[p] = addr
	}
	gw.rpLookup[addr] = rp
	return nil
}

func (gw *gateway) handleFunc(pattern string) (string, http.HandlerFunc) {
	return pattern, func(w http.ResponseWriter, r *http.Request) {
		span, ctx := tracer.StartSpanFromContext(r.Context(), ddconfig.GetService(), tracer.ResourceName(pattern))
		fields := []zap.Field{
			zap.Uint64("dd.trace_id", span.Context().TraceID()),
			zap.Uint64("dd.span_id", span.Context().SpanID()),
		}
		defer span.Finish()
		r = r.WithContext(ctx)
		if err := tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(r.Header)); err != nil {
			gw.l.Error("failed to inject span", append(fields, zap.Error(err))...)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to inject span: %v", err)))
			return
		}

		addr, ok := gw.addrLookup[pattern]
		if !ok {
			gw.l.Error("requested pattern is not found", append(fields, zap.String("pattern", pattern))...)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("%s is not found", pattern)))
			return
		}

		ip := httplib.GetIP(r)
		if ip == nil {
			gw.l.Error(fmt.Sprintf("ip address: %s is invalid format", ip.String()), fields...)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ip address is missing"))
			return
		}

		s, err := gw.ratelimit(w, r, pattern, ip)
		if err != nil {
			gw.l.Error("request failed", append(fields, zap.Error(err))...)
			w.WriteHeader(s)
			w.Write([]byte(err.Error()))
			return
		}

		gw.l.Info(fmt.Sprintf("reverse proxy: %s to %s:%s", pattern, addr, pattern), fields...)
		gw.rpLookup[addr].ServeHTTP(w, r)
	}
}

func (gw *gateway) ratelimit(w http.ResponseWriter, r *http.Request, pattern string, ip net.IP) (int, error) {
	sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to extract span context: %v", err)
	}
	span := tracer.StartSpan(ddconfig.GetService(ddconfig.WithServiceSuffix(".ratelimit")), tracer.ChildOf(sctx))
	defer span.Finish()

	key := fmt.Sprintf("%s %s", ip.String(), pattern)
	res, err := gw.limiter.Allow(r.Context(), key, getLimit())
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to allow request to %s from %s: %v", ip.String(), pattern, err)
	}
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
	if res.Allowed == 0 {
		retryAfter := strconv.Itoa(int(res.RetryAfter / time.Second))
		w.Header().Set("X-RateLimit-Reset", retryAfter)
		return http.StatusTooManyRequests, fmt.Errorf("exceeds rate limit, retry in %s second(s", retryAfter)
	}
	return http.StatusOK, nil
}

// For dynamic rate limit configuration.
var (
	defaultLimit    redisrate.Limit                      = redisrate.PerHour(60)
	limitByTimeUnit map[string]func(int) redisrate.Limit = map[string]func(int) redisrate.Limit{
		"second": redisrate.PerSecond,
		"Second": redisrate.PerSecond,
		"SECOND": redisrate.PerSecond,
		"minute": redisrate.PerMinute,
		"Minute": redisrate.PerMinute,
		"MINUTE": redisrate.PerMinute,
		"hour":   redisrate.PerHour,
		"Hour":   redisrate.PerHour,
		"HOUR":   redisrate.PerHour,
	}
)

func getLimit() redisrate.Limit {
	var unit string
	var limit string
	if unit = os.Getenv("RATELIMIT_TIME_UNIT"); unit == "" {
		return defaultLimit
	}
	if limit = os.Getenv("RATELIMIT_LIMIT"); limit == "" {
		return defaultLimit
	}
	l, ok := limitByTimeUnit[unit]
	if !ok {
		return defaultLimit
	}
	v, err := strconv.Atoi(limit)
	if err != nil {
		return defaultLimit
	}
	return l(v)
}
