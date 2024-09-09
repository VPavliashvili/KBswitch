package main

import (
	"fmt"
	"net/http"

	"kbswitch/internal/app"
	"kbswitch/internal/app/api"
)

// this is provided from build args
var compileDate string

func main() {
	a := app.New(compileDate)

	fmt.Printf("APPLICATION STARTED\n")

	router := api.InitRouter(a)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.Port),
		Handler: router,
	}

	server.ListenAndServe()
}
