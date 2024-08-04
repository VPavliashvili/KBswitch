package main

import (
	"fmt"
	"net/http"
	"time"

	"kbswitch/internal/app"
	"kbswitch/internal/app/api"
)

var compileDate string

func main() {
    compTime, err := time.Parse(time.RFC3339, compileDate)
    if err != nil {
        panic("could not parse build time\n" + err.Error())
    }

    // this will be separated out to the config
	app := app.Application{
		Config: app.Config{
			Port:         6012,
			ReadTimeout:  5,
			WriteTimeout: 5,
		},
        BuildDate: compTime,
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
