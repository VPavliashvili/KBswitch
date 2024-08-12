package api

import (
	"fmt"
	"net/http"

	"kbswitch/docs"
	"kbswitch/internal/app"
	"kbswitch/internal/app/api/controllers/switches"
	"kbswitch/internal/app/api/controllers/system"
	"kbswitch/internal/app/api/middlewares"
	"kbswitch/internal/app/api/router"

	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter(app app.Application) *router.CustomMux {

	docs.SwaggerInfo.Title = "Keyboard switches registry API"
	docs.SwaggerInfo.Description = "This is a backend of upcoming website"
	docs.SwaggerInfo.Version = "1.0"

	router := router.CreateAndSetup(func(this *router.CustomMux) *router.CustomMux {
		this.Use(middlewares.ContentTypeJSON)

		this.AddGroup("/api/system/", func(ng *router.Group) {
			c := system.New(app.BuildDate)

			ng.HandleRouteFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
				c.HandleAbout(w, r)
			})
		})

		this.AddGroup("/api/switches/", func(ng *router.Group) {
			c := switches.New(app.InjectedRepos.SwitchesRepo)

			ng.HandleRouteFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				c.HandleSwitches(w, r)
			})

			ng.HandleRouteFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
				c.HandleSwitchByID(w, r)
			})
		})

		this.HandleFunc("GET /swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", app.Config.Port)),
		))

		return this
	})

	return router
}
