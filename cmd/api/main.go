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
	compTime, err := time.Parse(time.RFC3339, compileDate)
	if err != nil {
		panic("could not parse build time\n" + err.Error())
	}

	// config part will be separated out to the config
	app := app.Application{
		Config: app.Config{
			Port:         6012,
			ReadTimeout:  5,
			WriteTimeout: 5,
		},
		BuildDate: compTime,
		Repos: app.InjectedRepos{
			Switches: nil, // after writing real implementation gotta create the instance here
		},
		Services: app.InjectedServices{
			Switches: nil, // after writing real implementation gotta create the instance here
		},
	}

	fmt.Printf("APPLICATION STARTED\n")

	router := api.InitRouter(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(app.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.WriteTimeout) * time.Second,
	}

	server.ListenAndServe()
}
