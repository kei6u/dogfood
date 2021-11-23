package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	redisrate "github.com/go-redis/redis_rate/v9"
	"github.com/kei6u/dogfood/driver"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

const (
	createRecordEndpoint = "/v1/dogfood/record"
	listRecordsEndpoint  = "/v1/dogfood/records"
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

	var addr string
	var dogfoodBackendAddr string

	if addr = os.Getenv("ADDR"); addr == "" {
		logger.Fatal("ADDR is missing")
	}
	if dogfoodBackendAddr = os.Getenv("DOGFOOD_BACKEND_ADDR"); dogfoodBackendAddr == "" {
		logger.Fatal("DOGFOOD_BACKEND_ADDR is missing")
	}

	ctx := context.Background()

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

	http.HandleFunc(createRecordEndpoint, gw.rateLimit(createRecordEndpoint))
	http.HandleFunc(listRecordsEndpoint, gw.rateLimit(listRecordsEndpoint))

	logger.Info("gateway has started", zap.String("port", addr))
	http.ListenAndServe(fmt.Sprintf(":%s", addr), nil)
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

// For dynamic rate limit configuration.
var (
	defaultLimit redisrate.Limit                      = redisrate.PerHour(60)
	limitByUnit  map[string]func(int) redisrate.Limit = map[string]func(int) redisrate.Limit{
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
	var vol string
	if unit = os.Getenv("RATELIMIT_UNIT"); unit == "" {
		return defaultLimit
	}
	if vol = os.Getenv("RATELIMIT_VOLUME"); vol == "" {
		return defaultLimit
	}
	l, ok := limitByUnit[unit]
	if !ok {
		return defaultLimit
	}
	v, err := strconv.Atoi(vol)
	if err != nil {
		return defaultLimit
	}
	return l(v)
}

func getIP(r *http.Request) net.IP {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		for _, h := range []string{"X-Forwarded-For", "x-forwarded-for", "X-FORWARDED-FOR"} {
			if ip = r.Header.Get(h); ip != "" {
				break
			}
		}
	}
	return net.ParseIP(ip)
}

func (gw *gateway) rateLimit(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipaddr := getIP(r)
		if ipaddr == nil {
			gw.l.Error(fmt.Sprintf("ipaddress: %s is invalid format", ipaddr))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ip address is missing"))
		}
		key := fmt.Sprintf("%s %s", ipaddr, pattern)
		res, err := gw.limiter.Allow(r.Context(), key, getLimit())
		if err != nil {
			gw.l.Error("failed to allow request", zap.Error(err), zap.String("IPAddress", ipaddr.String()), zap.String("pattern", pattern))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to allow request to %s from %s: %s", ipaddr.String(), pattern, err)))
			return
		}
		fields := []zap.Field{zap.String("RateLimit-Remaining", strconv.Itoa(res.Remaining))}
		w.Header().Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))
		if res.Allowed == 0 {
			retryAfter := strconv.Itoa(int(res.RetryAfter / time.Second))
			gw.l.Error("rate limited", zap.String("retryafter", retryAfter))
			w.Header().Set("RateLimit-RetryAfter-Second", retryAfter)
			fields = append(fields, zap.String("RateLimit-RetryAfter-Second", retryAfter))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf("retry after %s second(s)", retryAfter)))
			return
		}
		addr, ok := gw.addrLookup[pattern]
		if !ok {
			gw.l.Error("requested pattern is not found", zap.String("pattern", pattern))
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("%s is not found", pattern)))
			return
		}
		gw.l.Info(fmt.Sprintf("%s to %s:%s", pattern, addr, pattern), fields...)
		gw.rpLookup[addr].ServeHTTP(w, r)
	}
}
