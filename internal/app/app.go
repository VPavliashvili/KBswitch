package app

import (
	"kbswitch/internal/core/switches"
	"time"
)

type Application struct {
	Config    Config
	BuildDate time.Time
	Repos     InjectedRepos
	Services  InjectedServices
}

type InjectedRepos struct {
	Switches switches.Repo
}

type InjectedServices struct {
	Switches switches.Service
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type DbConfig struct {
	User string
	Pass string
	Host string
	Db   string
	Port uint16
}
