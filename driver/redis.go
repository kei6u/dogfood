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
	var pw string

	if host = os.Getenv("REDIS_HOST"); host == "" {
		return nil, nil, errors.New("redis host is missing")
	}
	if addr = os.Getenv("REDIS_ADDR"); addr == "" {
		return nil, nil, errors.New("redis addr is missing")
	}
	if pw = os.Getenv("REDIS_PASSWORD"); pw == "" {
		return nil, nil, errors.New("redis password is missing")
	}

	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, addr),
		Password: pw,
		DB:       0,
	})
	return c, c.Close, nil
}
