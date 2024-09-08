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
	"kbswitch/internal/app/database"
	switchservice "kbswitch/internal/pkg/switches"
	"kbswitch/internal/pkg/switches/repo"

	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter(app app.Application) *router.CustomMux {

	docs.SwaggerInfo.Title = "Keyboard switches registry API"
	docs.SwaggerInfo.Description = "This is a backend of upcoming website"
	docs.SwaggerInfo.Version = "0.0.1"

	router := router.CreateAndSetup(func(this *router.CustomMux) *router.CustomMux {
		this.Use(middlewares.ContentTypeJSON)
		this.Use(middlewares.Timeout((app.Config.Timeout)))
		this.Use(middlewares.InitPgxPool("switches", app.DbConfig))

		this.AddGroup("/api/system/", func(ng *router.Group) {
			c := system.New(app.BuildDate)

			ng.HandleRouteFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
				c.HandleAbout(w, r)
			})
		})

		this.AddGroup("/api/switches/", func(ng *router.Group) {
			c := switches.New(app.Services.Switches)

			ng.HandleRouteFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				app.Repos.Switches = repo.New(r.Context(), database.Get("switches"), app.DbConfig)
				app.Services.Switches = switchservice.New(app.Repos.Switches)
				c = switches.New(app.Services.Switches)

				c.HandleSwitches(r.Context(), w, r)
			})

			ng.HandleRouteFunc("GET /{brand}/{name}", func(w http.ResponseWriter, r *http.Request) {
				c.HandleSingleSwitch(r.Context(), w, r)
			})

			ng.HandleRouteFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
				c.HandleSwitchAdd(r.Context(), w, r)
			})

			ng.HandleRouteFunc("DELETE /{brand}/{name}", func(w http.ResponseWriter, r *http.Request) {
				c.HandleSwitchRemove(r.Context(), w, r)
			})

			ng.HandleRouteFunc("PATCH /{brand}/{name}", func(w http.ResponseWriter, r *http.Request) {
				c.HandleSwitchUpdate(r.Context(), w, r)
			})
		})

		this.HandleFunc("GET /swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", app.Config.Port)),
		))

		return this
	})

	return router
}
