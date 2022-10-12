package config

import (
	"fmt"

	"github.com/babyname/fate/ent"
)

const mysqlDSN = "mysql://%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true"
const sqlite3DSN = "file:%v?cache=shared&_journal=WAL&_fk=1"

type BuildFunc func(*database) (*ent.Client, error)

type Database interface {
	BuildClient() (*ent.Client, error)
}

var driverDSN = map[string]BuildFunc{
	"sqlite3": buildSqlite3,
	"mysql":   buildMysql,
}

func buildSqlite3(database *database) (*ent.Client, error) {
	dsn := sqlite3DSN
	if database.dsn != "" {
		dsn = database.dsn
	}
	return ent.Open(database.driver, fmt.Sprintf(dsn, database.DBName))
}

func buildMysql(database *database) (*ent.Client, error) {
	dsn := mysqlDSN
	if database.dsn != "" {
		dsn = database.dsn
	}
	return ent.Open(database.driver, fmt.Sprintf(dsn, database.User, database.Pwd, database.Host, database.Port, database.DBName))
}

type database struct {
	driver string `json:"driver"`
	dsn    string `json:"dsn,omitempty" json-getter:"DSN" json-setter:"SetDSN"`
	host   string `json:"host,omitempty"`
	port   string `json:"port,omitempty"`
	user   string `json:"user,omitempty"`
	pwd    string `json:"pwd,omitempty"`
	dbName string `json:"dbname,omitempty" json-getter:"DBName" json-setter:"SetDBName"`
}

func (d *database) Driver() string {
	return d.driver
}

func (d *database) SetDriver(driver string) {
	d.driver = driver
}

func (d *database) DSN() string {
	return d.dsn
}

func (d *database) SetDSN(dSN string) {
	d.dsn = dSN
}

func (d *database) Host() string {
	return d.host
}

func (d *database) SetHost(host string) {
	d.host = host
}

func (d *database) Port() string {
	return d.port
}

func (d *database) SetPort(port string) {
	d.port = port
}

func (d *database) User() string {
	return d.user
}

func (d *database) SetUser(user string) {
	d.user = user
}

func (d *database) Pwd() string {
	return d.pwd
}

func (d *database) SetPwd(pwd string) {
	d.pwd = pwd
}

func (d *database) DBName() string {
	return d.dbName
}

func (d *database) SetDBName(db_name string) {
	d.dbName = db_name
}

func (d *database) BuildClient() (*ent.Client, error) {
	fn, ok := driverDSN[d.driver]
	if !ok {
		return nil, fmt.Errorf("the driver of %v is not supported", d.driver)
	}
	return fn(d)
}
