package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kei6u/dogfood/driver"
	grpcbackend "github.com/kei6u/dogfood/grpc/backend"
	"go.uber.org/zap"
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

	s, err := grpcbackend.NewServer(
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
