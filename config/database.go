package config

type DBConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
	Log    string `json:"log"`
}
