package main

import (
	"fmt"
	"net/http"
	"time"

	"kbswitch/internal/app/api"
)

func main() {
	fmt.Printf("APPLICATION STARTED\n")

	router := api.InitRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 6012),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server.ListenAndServe()
}
