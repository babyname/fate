package config

type Database struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}
