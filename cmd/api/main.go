package main

import (
	"fmt"
	"net/http"
	"time"

	"kbswitch/internal/app"
	"kbswitch/internal/app/api"
)

// this is provided from build args
var compileDate string

func main() {
	bd, err := time.Parse(time.RFC3339, compileDate)
	if err != nil {
		panic("could not parse build time\n" + err.Error())
	}

	// config part will be separated out to the config
	a := app.Application{
		Config: app.Config{
			Port:    6012,
			Timeout: 10,
		},
		DbConfig: app.DbConfig{
			User: "admin",
			Pass: "test",
			Host: "database_switches",
			Db:   "switches_store",
			Port: 5432,
		},
		BuildDate: bd,
	}
	// app.App = a

	fmt.Printf("APPLICATION STARTED\n")

	router := api.InitRouter(a)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.Port),
		Handler: router,
	}

	server.ListenAndServe()
}
