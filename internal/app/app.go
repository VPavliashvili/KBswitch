package app

import (
	"kbswitch/internal/core/switches"
	"time"
)

type Application struct {
	Config           Config
	BuildDate        time.Time
	InjectedRepos    InjectedRepos
	InjectedServices InjectedServices
}

type InjectedRepos struct {
	SwitchesRepo switches.Repo
}

type InjectedServices struct {
	SwitchesService switches.Service
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}
