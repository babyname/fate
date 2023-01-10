package config

const (
	defaultLogName = "fate.log"
	defaultLogPath = "log"
)

type LogConfig struct {
	Path       string `json:"path,omitempty"`
	LogType    string `json:"log_type,omitempty"`
	ShowSource bool   `json:"show_source,omitempty"`
	Level      string `json:"level,omitempty"`
}

func defaultLogConfig() LogConfig {
	return LogConfig{
		Path:       "log/fate.log",
		ShowSource: false,
		Level:      "INFO",
		LogType:    "json",
	}
}
