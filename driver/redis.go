package driver

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// NewRedis returns a client to the Redis Server specified by Options.
func NewRedis(ctx context.Context) (*redis.Client, func() error, error) {
	var host string
	var addr string

	if host = os.Getenv("REDIS_HOST"); host == "" {
		return nil, nil, errors.New("host is missing")
	}
	if addr = os.Getenv("REDIS_ADDR"); addr == "" {
		return nil, nil, errors.New("addr is missing")
	}

	c := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, addr),
		DB:   0,
	})
	return c, c.Close, nil
}
