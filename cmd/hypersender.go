package main

import (
	"context"
	"github.com/drblez/hypersender/config"
	"github.com/drblez/hypersender/logger"
	"github.com/drblez/hypersender/worker"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func buildContainer() *dig.Container {
	c := dig.New()
	_ = c.Provide(config.Init)
	_ = c.Provide(logger.Init)
	_ = c.Provide(worker.Init)
	_ = c.Provide(worker.InitFSWorker)
	return c
}

func do(config *config.Config, log *logrus.Entry, worker *worker.Worker) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	err := worker.Do(ctx, wg)
	if err != nil {
		panic(err)
	}
	<-quit
	cancel()
	wg.Wait()
}

func main() {
	c := buildContainer()
	err := c.Invoke(do)
	if err != nil {
		panic(err)
	}
}
