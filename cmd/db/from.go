package main

const dsn = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

type From struct {
	DSN          string `json:"dsn"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	DBName       string `json:"dbname"`
	MaxIdleCon   int    `json:"max_idle_con"`
	MaxOpenCon   int    `json:"max_open_con"`
	Driver       string `json:"driver"`
	File         string `json:"file"`
	ShowSQL      bool   `json:"show_sql"`
	ShowExecTime bool   `json:"show_exec_time"`
}
