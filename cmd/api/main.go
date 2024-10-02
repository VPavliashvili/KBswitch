package main

import (
	"fmt"
	"net/http"

	"kbswitch/internal/app"
	"kbswitch/internal/app/api"
	"kbswitch/internal/core/common/logger"
)

// this is provided from build args
var compileDate string

func main() {
	a := app.New(compileDate)
	logger.Init(a)

	logger.Info("APPLICATION STARTED")

	router := api.InitRouter(a)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.Port),
		Handler: router,
	}

	server.ListenAndServe()
}
