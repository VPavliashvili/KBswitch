package app

import (
	"os"
	"strconv"
	"time"
)

const (
	APP_TIMEOUT = "APP_TIMEOUT"
	APP_PORT    = "APP_PORT"
	APP_DB_USER = "APP_DB_USER"
	APP_DB_PASS = "APP_DB_PASS"
	APP_DB_HOST = "APP_DB_HOST"
	APP_DB_PORT = "APP_DB_PORT"
	APP_DB      = "APP_DB"
)

type Application struct {
	Config    Config
	DbConfig  DbConfig
	BuildDate time.Time
}

type Config struct {
	Timeout int
	Port    int
}

type DbConfig struct {
	User string
	Pass string
	Host string
	Db   string
	Port int
}

func New(buildDate string) Application {
	timeout, _ := strconv.Atoi(os.Getenv(APP_TIMEOUT))
	port, _ := strconv.Atoi(os.Getenv(APP_PORT))
	user := os.Getenv(APP_DB_USER)
	pass := os.Getenv(APP_DB_PASS)
	host := os.Getenv(APP_DB_HOST)
	db := os.Getenv(APP_DB)
	dbp, _ := strconv.Atoi(os.Getenv(APP_DB_PORT))

	bd, err := time.Parse(time.RFC3339, buildDate)
	if err != nil {
		panic("could not parse build time\n" + err.Error())
	}

	return Application{
		Config: Config{
			Timeout: timeout,
			Port:    port,
		},
		DbConfig: DbConfig{
			User: user,
			Pass: pass,
			Host: host,
			Db:   db,
			Port: dbp,
		},
		BuildDate: bd,
	}
}
