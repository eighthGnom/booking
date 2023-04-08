package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/eighthGnom/booking/internal/config"
	"github.com/eighthGnom/booking/internal/handlers"
	"github.com/eighthGnom/booking/internal/models"
	"github.com/eighthGnom/booking/internal/render"
)

var port = ":8080"
var appConfig config.AppConfig
var sessionManager *scs.SessionManager

// main is the main application function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("Starting app on port %s", port))
	srv := &http.Server{
		Addr:    port,
		Handler: routes(&appConfig),
	}
	log.Fatal(srv.ListenAndServe())

}

func run() error {
	var err error

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	// change this to true when in production
	appConfig.InProduction = false

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = appConfig.InProduction

	appConfig.SessionManager = sessionManager
	appConfig.TemplateCache, err = render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	appConfig.UseCache = false
	repo := handlers.NewRepo(&appConfig)
	handlers.SetRepoForHandlers(repo)
	render.SetTemplatesConfig(&appConfig)

	return nil
}
