package logger

import (
	"context"
	"github.com/drblez/hypersender/config"
	"github.com/drblez/hypersender/utils"
	"github.com/drblez/tasks"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

func Init(config *config.Config) *logrus.Entry {
	log := logrus.New()

	path := config.LogPath + "/hypersender.log"

	if err := utils.MakeDirAll(path); err != nil {
		panic(err)
	}

	writer, err := rotatelogs.New(path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		panic(err)
	}

	log.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},
		&logrus.JSONFormatter{},
	))

	if config.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	if config.Console {
		log.Formatter = &logrus.TextFormatter{}
	}
	return log.WithContext(context.Background())
}

type TasksLogger struct {
	*logrus.Entry
}

func NewTaskLogger(log *logrus.Entry) *TasksLogger {
	tl := &TasksLogger{}
	tl.Entry = log.WithContext(context.Background())
	return tl
}

func (logger *TasksLogger) AddField(key string, value interface{}) tasks.Logger {
	tl := &TasksLogger{}
	tl.Entry = logger.WithField(key, value)
	return tl
}
