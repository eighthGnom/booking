package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/eighthGnom/booking/internal/config"
	"github.com/eighthGnom/booking/internal/models"
	"github.com/justinas/nosurf"
)

var appConfig *config.AppConfig
var functions = template.FuncMap{}

// SetTemplatesConfig sets the config for the template package
func SetTemplatesConfig(config *config.AppConfig) {
	appConfig = config
}

// AddDefaultData adds to the templates default data
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = appConfig.SessionManager.PopString(r.Context(), "flash")
	td.Error = appConfig.SessionManager.PopString(r.Context(), "error")
	td.Warning = appConfig.SessionManager.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate render templates and writes them to response
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	tc := map[string]*template.Template{}
	if appConfig.UseCache {
		tc = appConfig.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	td = AddDefaultData(td, r)

	err := t.Execute(w, td)
	if err != nil {
		fmt.Println("error while render template", err)
		return
	}

}

// CreateTemplateCache prepare templates and adds it to the cash
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCash := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCash, err
	}
first:
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCash, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCash, err
		}
		if len(matches) != 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCash, err
			}
		} else {
			continue first
		}

		myCash[name] = ts
	}
	return myCash, nil

}
