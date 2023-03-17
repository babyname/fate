package config

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

func defaultDBSqlite3() DBConfig {
	return DBConfig{
		Name:   "fate",
		Driver: "sqlite3",
	}
}

func defaultDBMysql() DBConfig {
	return DBConfig{
		Host:   "127.0.0.1",
		Port:   "3306",
		User:   "root",
		Pwd:    "111111",
		Name:   "fate",
		Driver: "mysql",
	}
}
