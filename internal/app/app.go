package app

import (
	"time"
)

// state of the application, is singleton
// var App Application

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
