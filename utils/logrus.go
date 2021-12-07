package utils

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type MyHook struct {
}

type Data struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func InitLogrus() {

	path := "routinelog"

	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationCount(5),
		rotatelogs.WithRotationTime(time.Duration(76400)*time.Second),
	)

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	writers := []io.Writer{
		writer,
		os.Stdout,
	}

	fileAndStdoutWriter := io.MultiWriter(writers...)

	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	} else {
		logrus.Info("failed to log to file.")
	}

	// logrus.SetLevel(logrus.InfoLevel)
	logrus.SetLevel(logrus.WarnLevel)
	// logrus.SetLevel(logrus.ErrorLevel)
	// logrus.SetLevel(logrus.InfoLevel)
}

func (h *MyHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
