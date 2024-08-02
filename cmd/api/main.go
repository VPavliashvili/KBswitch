package main

import (
	"fmt"
	"net/http"
	"time"

	"vpavliashvili.mech-switch/internal/app/api/controllers/system"
	"vpavliashvili.mech-switch/internal/app/api/middlewares"
	"vpavliashvili.mech-switch/internal/app/api/router"
)

func main() {
	fmt.Printf("APPLICATION STARTED\n")

	router := router.CreateAndSetup(func(this *router.CustomMux) *router.CustomMux {
		this.Use(middlewares.ContentTypeJSON)

		this.AddGroup("/api/system/", func(ng *router.Group) {
			c := system.New()

			ng.HandleRouteFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
				c.HandleAbout(w, r)
			})
		})

		return this

	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 6012),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server.ListenAndServe()
}
