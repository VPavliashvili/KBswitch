package api

import (
	"net/http"

	"kbswitch/internal/app/api/controllers/system"
	"kbswitch/internal/app/api/middlewares"
	"kbswitch/internal/app/api/router"
)

func InitRouter() *router.CustomMux {
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

	return router
}
