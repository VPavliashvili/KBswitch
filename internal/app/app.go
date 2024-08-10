package app

import (
	"kbswitch/internal/core/repositories"
	"time"
)

type Application struct {
	Config        Config
	BuildDate     time.Time
	InjectedRepos InjectedRepos
}

type InjectedRepos struct {
	SwitchesRepo repositories.SwitchesRepo
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}
