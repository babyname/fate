package fate

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
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
	head []string
	path string
	file *os.File
}

func (l *logInformation) Finish() error {
	return l.sugar.Sync()
}

func initOutputWithConfig(output config.FileOutput) Information {
	switch output.OutputMode {
	//case config.OutputModelJSON:
	//	return jsonOutput(output.Path)
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
		out := headNameJSONOutput(j.head, n, nil)
		//output json
		_, _ = w.Write(out)
		_, _ = w.WriteString(",\n")
	}
	return w.Flush()

}
func (j *jsonInformation) Head(heads ...string) error {
	j.head = heads
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
func (l *logInformation) Write(names ...Name) error {
	for _, n := range names {
		out := headNameOutput(l.head, n, func(s string) bool {
			return s == "姓名"
		})
		l.sugar.Infow(n.String(), out...)
	}
	return nil
}

func (l *logInformation) Head(heads ...string) error {
	l.head = heads
	return nil
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
		out := nameOutputString(c.head, n)
		_ = w.Write(out)
	}
	w.Flush()
	return nil
}

func (c *csvInformation) Finish() error {
	return c.file.Close()
}

func (c *csvInformation) Head(heads ...string) (e error) {
	c.head = heads
	w := csv.NewWriter(c.file)
	e = w.Write(heads)
	if e != nil {
		return e
	}
	w.Flush()
	return nil
}

func headNameOutput(heads []string, name Name, skip func(string) bool) (out []interface{}) {
	for _, h := range heads {
		if skip != nil && skip(h) {
			continue
		}
		switch h {
		case "姓名":
			out = append(out, h, name.String())
		case "笔画":
			out = append(out, h, name.Strokes())
		case "拼音":
			out = append(out, h, name.PinYin())
		case "喜用神":
			out = append(out, h, name.XiYongShen())
		case "八字", "生辰八字":
			out = append(out, h, name.BaZi())
		}
	}
	return
}

func headNameJSONOutput(heads []string, name Name, skip func(string) bool) (b []byte) {
	out := make(map[string]string)
	for _, h := range heads {
		if skip != nil && skip(h) {
			continue
		}
		switch h {
		case "姓名":
			out[h] = name.String()
		case "笔画":
			out[h] = name.Strokes()
		case "拼音":
			out[h] = name.PinYin()
		case "喜用神":
			out[h] = name.XiYongShen()
		}
	}
	by, e := json.Marshal(out)
	if e != nil {
		return nil
	}
	return by
}

func headNameOutputString(heads []string, name Name, skip func(string) bool) (out []string) {
	for _, h := range heads {
		if skip != nil && skip(h) {
			continue
		}
		switch h {
		case "姓名":
			out = append(out, h, name.String())
		case "笔画":
			out = append(out, h, name.Strokes())
		case "拼音":
			out = append(out, h, name.PinYin())
		case "喜用神":
			out = append(out, h, name.XiYongShen())
		}
	}
	return
}
func nameOutputString(heads []string, name Name) (out []string) {
	for _, h := range heads {
		switch h {
		case "姓名":
			out = append(out, name.String())
		case "笔画":
			out = append(out, name.Strokes())
		case "拼音":
			out = append(out, name.PinYin())
		case "喜用神":
			out = append(out, name.XiYongShen())
		}
	}
	return
}
