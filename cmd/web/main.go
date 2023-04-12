package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/eighthGnom/booking/internal/config"
	"github.com/eighthGnom/booking/internal/driver"
	"github.com/eighthGnom/booking/internal/handlers"
	"github.com/eighthGnom/booking/internal/helpers"
	"github.com/eighthGnom/booking/internal/models"
	"github.com/eighthGnom/booking/internal/render"
)

var port = ":8080"
var appConfig config.AppConfig
var sessionManager *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	fmt.Println(fmt.Sprintf("Starting app on port %s", port))
	srv := &http.Server{
		Addr:    port,
		Handler: routes(&appConfig),
	}
	log.Fatal(srv.ListenAndServe())

}

func run() (*driver.DB, error) {
	var err error

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")
	if err != nil {
		return nil, err
	}

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	// change this to true when in production
	appConfig.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

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
	helpers.SetHelpersConfig(&appConfig)

	return db, nil
}
