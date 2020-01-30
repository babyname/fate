package fate

import (
	"bufio"
	"fmt"
	"github.com/godcong/fate/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Information interface {
	Write(names ...Name) error
	Finish() error
}

type jsonInformation struct {
	path string
	file *os.File
}

type logInformation struct {
	path  string
	sugar *zap.SugaredLogger
}

func (l *logInformation) Write(names ...Name) error {
	for _, n := range names {
		l.sugar.Infow("Information-->", "名字", n.String(), "笔画", n.Strokes(), "拼音", n.PinYin())
	}
	return nil
}

func (l *logInformation) Finish() error {
	return l.sugar.Sync()
}

func NewWithConfig(cfg config.Config) Information {
	switch cfg.FileOutput.OutputMode {
	case config.OutputModelJSON:
		return jsonOutput(cfg.FileOutput.Path)
	}

	return logOutput(cfg.FileOutput.Path)
}

func (j *jsonInformation) Finish() error {
	return j.file.Close()
}

func (j *jsonInformation) Write(names ...Name) error {
	w := bufio.NewWriter(j.file)
	for _, n := range names {
		_, _ = w.WriteString(n.String())
		_, _ = w.WriteString("\r\n")
	}
	return w.Flush()

}

func jsonOutput(path string) Information {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0755)
	if e != nil {
		panic(fmt.Errorf("json output failed:%w", e))
	}
	return &jsonInformation{
		path: path,
		file: file,
	}
}

func logOutput(path string) Information {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "",
		TimeKey:        "",
		NameKey:        "",
		CallerKey:      "",
		StacktraceKey:  "",
		LineEnding:     "",
		EncodeLevel:    nil,
		EncodeTime:     nil,
		EncodeDuration: nil,
		EncodeCaller:   nil,
		EncodeName:     nil,
	}
	cfg.OutputPaths = []string{
		path,
	}
	cfg.ErrorOutputPaths = []string{
		"out.log",
	}

	cfg.DisableCaller = true
	cfg.DisableStacktrace = true

	logger, e := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if e != nil {
		panic(e)
	}
	return &logInformation{
		path:  path,
		sugar: logger.Sugar(),
	}
}
