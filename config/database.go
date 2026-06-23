package config

const (
	InitModeAuto = "auto"
	InitModeDB   = "db"
	InitModeJSON = "json"
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

	InitMode string `json:"init_mode,omitempty" yaml:"init_mode"`
	DBFile   string `json:"db_file,omitempty"   yaml:"db_file"`
}
