package app

import (
	"kbswitch/internal/core/switches"
	"time"
)

// state of the application, is singleton
// var App Application

type Application struct {
	Config    Config
	DbConfig  DbConfig
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
