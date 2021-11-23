package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	redisrate "github.com/go-redis/redis_rate/v9"
	"github.com/kei6u/dogfood/pkg/httplib"
	"go.uber.org/zap"
)

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
		addr, ok := gw.addrLookup[pattern]
		if !ok {
			gw.l.Error("requested pattern is not found", zap.String("pattern", pattern))
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("%s is not found", pattern)))
			return
		}

		ip := httplib.GetIP(r)
		if ip == nil {
			gw.l.Error(fmt.Sprintf("ip address: %s is invalid format", ip.String()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ip address is missing"))
			return
		}

		s, err := gw.ratelimit(w, r, pattern, ip)
		if err != nil {
			gw.l.Error("request failed", zap.Error(err))
			w.WriteHeader(s)
			w.Write([]byte(err.Error()))
			return
		}

		gw.l.Info(fmt.Sprintf("%s to %s:%s", pattern, addr, pattern))
		gw.rpLookup[addr].ServeHTTP(w, r)
	}
}

func (gw *gateway) ratelimit(w http.ResponseWriter, r *http.Request, pattern string, ip net.IP) (int, error) {
	key := fmt.Sprintf("%s %s", ip.String(), pattern)
	res, err := gw.limiter.Allow(r.Context(), key, getLimit())
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to allow request to %s from %s: %v", ip.String(), pattern, err)
	}
	w.Header().Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))
	if res.Allowed == 0 {
		retryAfter := strconv.Itoa(int(res.RetryAfter / time.Second))
		w.Header().Set("RateLimit-RetryAfter-Second", retryAfter)
		return http.StatusTooManyRequests, fmt.Errorf("exceeds rate limit, retry in %s second(s", retryAfter)
	}
	return http.StatusOK, nil
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
