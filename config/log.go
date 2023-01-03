package config

type Logger struct {
	Path       string `json:"path,omitempty"`
	ShowSource bool   `json:"show_source,omitempty"`
	Level      string `json:"level,omitempty"`
}
