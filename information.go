package fate

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/godcong/fate/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Information interface {
	Write(names ...Name) error
	Head(heads ...string) error
	Finish() error
}

type jsonInformation struct {
	head []string
	path string
	file *os.File
}

type logInformation struct {
	path  string
	sugar *zap.SugaredLogger
	head  []string
}

type csvInformation struct {
	path string
	file *os.File
}

func (l *logInformation) Write(names ...Name) error {
	for _, n := range names {
		l.sugar.Infow(n.String(), "笔画", n.Strokes(), "拼音", n.PinYin())
	}
	return nil
}

func (l *logInformation) Finish() error {
	return l.sugar.Sync()
}

func initOutputWithConfig(output config.FileOutput) Information {
	switch output.OutputMode {
	case config.OutputModelJSON:
		return jsonOutput(output.Path)
	case config.OutputModeCSV:
		return csvOutput(output.Path)
	}

	return logOutput(output.Path)
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
func (j *jsonInformation) Head(heads ...string) error {
	j.head = heads
	return nil
}

func (l *logInformation) Head(heads ...string) error {
	l.head = heads
	return nil
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
		MessageKey:     "name",
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
	//cfg.ErrorOutputPaths = []string{
	//	"out.log",
	//}

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

func csvOutput(path string) Information {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0755)
	if e != nil {
		panic(fmt.Errorf("json output failed:%w", e))
	}

	return &csvInformation{
		path: path,
		file: file,
	}
}

func (c *csvInformation) Write(names ...Name) error {
	w := csv.NewWriter(c.file)
	for _, n := range names {
		_ = w.Write([]string{
			n.String(), n.Strokes(), n.PinYin(), n.WuXing(),
		})
	}
	w.Flush()
	return nil
}

func (c *csvInformation) Finish() error {
	return c.file.Close()
}

func (c *csvInformation) Head(heads ...string) (e error) {
	w := csv.NewWriter(c.file)
	e = w.Write(heads)
	if e != nil {
		return e
	}
	w.Flush()
	return nil
}
