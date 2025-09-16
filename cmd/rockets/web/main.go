package main

import (
	"context"
	"fmt"
	"lunar-backend-engineer-challenge/cmd/di"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	rocketsDI := di.Init()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	defer func() {
		rocketsDI.Services.Logger.Warn("stopping application")
		stop()
	}()

	di.RunMigrations(ctx, rocketsDI.Services)

	if err := rocketsDI.Services.Router.ListenAndServe(fmt.Sprintf("%s:%v", rocketsDI.Config.ServerHost, rocketsDI.Config.ServerPort)); err != nil {
		rocketsDI.Services.Logger.Fatal("Error starting server", zap.Error(err))
	}
}
