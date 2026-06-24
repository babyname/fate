package config

import (
	"os"
	"path/filepath"
)

const (
	InitModeAuto = "auto"
	InitModeDB   = "db"
	InitModeJSON = "json"
)

const (
	DefaultDBReleaseURL = "https://github.com/babyname/fate/releases/latest/download"
)

type DBConfig struct {
	Driver  string `json:"driver,omitempty"  yaml:"driver"`
	DSN     string `json:"dsn,omitempty"     yaml:"dsn"`
	Mode    string `json:"mode,omitempty"    yaml:"mode"`
	Host    string `json:"host,omitempty"    yaml:"host"`
	Port    string `json:"port,omitempty"    yaml:"port"`
	User    string `json:"user,omitempty"    yaml:"user"`
	Pwd     string `json:"pwd,omitempty"     yaml:"pwd"`
	Name    string `json:"name,omitempty"    yaml:"name"`
	Timeout int    `json:"timeout,omitempty" yaml:"timeout"`

	InitMode    string `json:"init_mode,omitempty"    yaml:"init_mode"`
	DBFile      string `json:"db_file,omitempty"      yaml:"db_file"`
	AutoDownload bool  `json:"auto_download,omitempty" yaml:"auto_download"`
	DownloadURL string `json:"download_url,omitempty" yaml:"download_url"`
}

func (c *DBConfig) GetDataDir() string {
	if c.DBFile != "" {
		return filepath.Dir(c.DBFile)
	}
	exePath, err := os.Executable()
	if err == nil {
		return filepath.Dir(exePath)
	}
	return "."
}

func (c *DBConfig) GetDBFile() string {
	if c.DBFile != "" {
		return c.DBFile
	}
	name := c.Name
	if name == "" {
		name = "fate"
	}
	return filepath.Join(c.GetDataDir(), name+".db")
}

func (c *DBConfig) GetDBGZFile() string {
	return c.GetDBFile() + ".gz"
}

func (c *DBConfig) GetDownloadURL() string {
	if c.DownloadURL != "" {
		return c.DownloadURL
	}
	if envURL := os.Getenv("FATE_DB_DOWNLOAD_URL"); envURL != "" {
		return envURL
	}
	return DefaultDBReleaseURL
}

func (c *DBConfig) ShouldAutoDownload() bool {
	if c.AutoDownload {
		return true
	}
	if envAuto := os.Getenv("FATE_DB_AUTO_DOWNLOAD"); envAuto == "1" || envAuto == "true" {
		return true
	}
	return c.InitMode == "" || c.InitMode == InitModeAuto
}
