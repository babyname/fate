package config

type Database struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Pwd        string `json:"pwd"`
	Name       string `json:"name"`
	MaxIdleCon int    `json:"max_idle_con"`
	MaxOpenCon int    `json:"max_open_con"`
	Driver     string `json:"driver"`
	File       string `json:"file"`
	Dsn        string `json:"dsn"`
}
