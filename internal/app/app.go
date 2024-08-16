package app

import (
	"kbswitch/internal/core/switches/repositories"
	"kbswitch/internal/core/switches/services"
	"time"
)

type Application struct {
	Config           Config
	BuildDate        time.Time
	InjectedRepos    InjectedRepos
	InjectedServices InjectedServices
}

type InjectedRepos struct {
	SwitchesRepo repositories.SwitchesRepo
}

type InjectedServices struct {
	SwitchesService services.SwitchesService
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}
