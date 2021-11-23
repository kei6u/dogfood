package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/kei6u/dogfood/driver"
	dogfoodgrpc "github.com/kei6u/dogfood/grpc"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

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

	db, closeDB, err := driver.NewPsql()
	if err != nil {
		logger.Warn("exit due to connection failure of database", zap.Error(err))
		return
	}
	defer closeDB()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		logger.Info("Server shut down in 5 seconds")
		time.Sleep(5 * time.Second)
		cancel()
	}()

	go func() {
		dogfoodgrpc.NewServer(os.Getenv("GRPC_ADDR"), os.Getenv("GRPC_GATEWAY_ADDR"), logger, db).Start(ctx)
	}()

	<-ctx.Done()
}
