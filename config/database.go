package config

type Database struct {
	Driver string `json:"driver"`
	Addr   string `json:"addr"`
	User   string `json:"user"`
	Pwd    string `json:"pwd"`
	Name   string `json:"name"`
	Loc    string `json:"loc"`
	File   string `json:"file"`
	Dsn    string `json:"dsn"`
	Debug  bool   `json:"debug"`
	Log    string `json:"log"`
}
