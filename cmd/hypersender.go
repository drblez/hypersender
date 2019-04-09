package main

import (
	"github.com/drblez/hypersender/config"
	"github.com/drblez/hypersender/logger"
	"github.com/drblez/hypersender/worker"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
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
	log.Debugf("Start")
	err := worker.Do()
	if err != nil {
		panic(err)
	}
	log.Debugf("Done")
}

func main() {
	c := buildContainer()
	err := c.Invoke(do)
	if err != nil {
		panic(err)
	}
}
