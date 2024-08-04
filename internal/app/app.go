package app

import "time"

type Application struct {
	Config    Config
	BuildDate time.Time
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}
