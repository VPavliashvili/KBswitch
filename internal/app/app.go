package app

import (
	"os"
	"strconv"
	"time"
)

const (
	APP_TIMEOUT        = "APP_TIMEOUT"
	APP_PORT           = "APP_PORT"
	APP_DB_USER        = "APP_DB_USER"
	APP_DB_PASS        = "APP_DB_PASS"
	APP_DB_HOST        = "APP_DB_HOST"
	APP_DB_PORT        = "APP_DB_PORT"
	APP_DB             = "APP_DB"
	LOG_PATH           = "LOG_PATH"
	LOG_ENABLE_CONSOLE = "LOG_ENABLE_CONSOLE"
)

type Application struct {
	Config    Config
	DbConfig  DbConfig
	BuildDate time.Time
	Logging   Logging
}

type Config struct {
	Timeout int
	Port    int
}

type Logging struct {
	EnableConsole bool
	LogFilePath   string
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

	logpath := os.Getenv(LOG_PATH)
	hasConsole, _ := strconv.ParseBool(os.Getenv(LOG_ENABLE_CONSOLE))

	bd, err := time.Parse(time.RFC3339, buildDate)
	if err != nil {
		panic("could not parse build time\n" + err.Error())
	}

	return Application{
		Config: Config{
			Timeout: timeout,
			Port:    port,
		},
		Logging: Logging{
			LogFilePath:   logpath,
			EnableConsole: hasConsole,
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
