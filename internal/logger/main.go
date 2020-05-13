package logger

import (
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var Log = logrus.New()

type Fields = logrus.Fields

func init() {

	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filePath := strings.Split(frame.File, "/")
			return "", fmt.Sprintf("%s/%s:%d", filePath[len(filePath)-2], filePath[len(filePath)-1], frame.Line)
		},
	})
	Log.SetReportCaller(true)
	Log.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(config.C.LogLevel)
	if err != nil {
		Log.Errorf("cannot parse logLevel %s", config.C.LogLevel)
		lvl = logrus.InfoLevel
	}

	Log.SetLevel(lvl)
}

func New(name string) *logrus.Entry {
	return Log.WithField("source", name)
}
