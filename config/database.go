package config

// DBConfig 数据库连接配置
type DBConfig struct {
	Driver  string `json:"driver,omitempty"`
	DSN     string `json:"dsn,omitempty"`
	Host    string `json:"host,omitempty"`
	Port    string `json:"port,omitempty"`
	User    string `json:"user,omitempty"`
	Pwd     string `json:"pwd,omitempty"`
	Name    string `json:"name,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}
