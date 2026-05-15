package config

// LogConfig 日志配置
type LogConfig struct {
	Path       string `json:"path,omitempty"`
	LogType    string `json:"log_type,omitempty"`
	ShowSource bool   `json:"show_source,omitempty"`
	Level      string `json:"level,omitempty"`
}
