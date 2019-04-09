package worker

import (
	"context"
	"fmt"
	"github.com/drblez/hypersender/config"
	"github.com/drblez/hypersender/logger"
	"github.com/drblez/tasks"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

type Worker struct {
	config     *config.Config
	log        *logrus.Entry
	dispatcher *tasks.Dispatcher
	nw         *netWorker
}

func Init(config *config.Config, log *logrus.Entry, nw *netWorker) *Worker {
	return &Worker{
		config:     config,
		log:        log,
		dispatcher: tasks.NewDispatcher(config.FSParallelism, config.FSParallelism, logger.NewTaskLogger(log)),
		nw:         nw,
	}
}

func (fsw *Worker) doDir(startPath string) error {
	allDirs, err := ioutil.ReadDir(startPath)
	if err != nil {
		return err
	}

	netFunc := func(fileName string, file io.Reader) func() error {
		return func() error {
			fsw.nw.log.Debugf("Send file %s to %s", fileName, fsw.nw.config.URL)
			result, err := http.Post(fsw.nw.config.URL, fsw.nw.config.ContentType, file)
			if err != nil {
				return err
			}
			if result.StatusCode != http.StatusOK {
				return fmt.Errorf("bad status code: %d", result.StatusCode)
			}
			return nil
		}
	}

	fsFunc := func(fileName string) func() error {
		return func() error {
			fsw.log.Debugf("Process file: %s", fileName)
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			fsw.nw.dispatcher.Payload(netFunc(fileName, file))
			return nil
		}
	}

	for _, item := range allDirs {
		itemName := path.Join(startPath, item.Name())
		if item.IsDir() {
			err := fsw.doDir(itemName)
			if err != nil {
				return err
			}
			continue
		}
		fsw.dispatcher.Payload(fsFunc(itemName))
	}
	return nil
}

func (fsw *Worker) Do(ctx context.Context, wg *sync.WaitGroup) error {
	netCtx, netCancel := context.WithCancel(ctx)
	netWg := &sync.WaitGroup{}
	fsw.nw.dispatcher.Run(netCtx, netWg)
	fsw.dispatcher.Run(ctx, wg)
	err := fsw.doDir(fsw.config.Path)
	if err != nil {
		return err
	}
	netCancel()
	netWg.Wait()
	return nil
}

type netWorker struct {
	config     *config.Config
	log        *logrus.Entry
	dispatcher *tasks.Dispatcher
}

func InitFSWorker(config *config.Config, log *logrus.Entry) *netWorker {
	return &netWorker{
		config:     config,
		log:        log,
		dispatcher: tasks.NewDispatcher(config.NetParallelism, config.NetParallelism, logger.NewTaskLogger(log)),
	}
}
