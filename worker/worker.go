package worker

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/drblez/hypersender/config"
	"github.com/drblez/hypersender/logger"
	"github.com/drblez/tasks"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

var (
	Errors     = errorx.NewNamespace("worker")
	FileErrors = Errors.NewType("file_error")
	NetErrors  = Errors.NewType("net_error")
)

type Worker struct {
	config     *config.Config
	log        *logrus.Entry
	dispatcher *tasks.Dispatcher
	nw         *netWorker
}

func Init(config *config.Config, log *logrus.Entry, nw *netWorker) *Worker {
	dispatcher := tasks.NewDispatcher(config.FSParallelism, config.FSParallelism, logger.NewTaskLogger(log))
	return &Worker{
		config:     config,
		log:        log,
		dispatcher: dispatcher,
		nw:         nw,
	}
}

func (fsw *Worker) doDir(startPath string) error {
	allDirs, err := ioutil.ReadDir(startPath)
	if err != nil {
		return err
	}

	netFunc := func(fileName string, file io.ReadCloser) func() error {
		return func() error {
			defer file.Close()
			fsw.nw.log.Infof("Sending file %s...", fileName)
			result, err := http.Post(fsw.nw.config.URL, fsw.nw.config.ContentType, file)
			if err != nil {
				err := NetErrors.WrapWithNoMessage(err)
				if fsw.config.PanicOnErrors {
					panic(err)
				}
				return err
			}
			if !fsw.config.IgnoreServiceErrors {
				if result.StatusCode != http.StatusOK {
					err := NetErrors.New("bad status code: %d", result.StatusCode)
					if fsw.config.PanicOnErrors {
						panic(err)
					}
					return err
				}
			}
			return nil
		}
	}

	fsFunc := func(fileName string) func() error {
		return func() error {
			fsw.log.Debugf("Process file: %s", fileName)
			file, err := os.Open(fileName)
			if err != nil {
				err := FileErrors.WrapWithNoMessage(err)
				if fsw.config.PanicOnErrors {
					panic(err)
				}
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
		fsw.log.Debugf("Sent to payload queue: %s", itemName)
		fsw.dispatcher.Payload(fsFunc(itemName))
	}

	return nil
}

func (fsw *Worker) Do() error {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	netCtx, netCancel := context.WithCancel(context.Background())
	netWg := &sync.WaitGroup{}
	fsw.nw.dispatcher.Run(netCtx, netWg)
	fsw.dispatcher.Run(ctx, wg)
	err := fsw.doDir(fsw.config.Path)
	if err != nil {
		return err
	}
	fsw.dispatcher.Payload(func() error {
		cancel()
		return nil
	})
	wg.Wait()
	fsw.nw.dispatcher.Payload(func() error {
		netCancel()
		return nil
	})
	netWg.Wait()
	return nil
}

type netWorker struct {
	config     *config.Config
	log        *logrus.Entry
	dispatcher *tasks.Dispatcher
}

func InitFSWorker(config *config.Config, log *logrus.Entry) *netWorker {
	dispatcher := tasks.NewDispatcher(config.NetParallelism, config.NetParallelism, logger.NewTaskLogger(log))
	//dispatcher.QuitOnEmpty()
	return &netWorker{
		config:     config,
		log:        log,
		dispatcher: dispatcher,
	}
}
