package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/eighthGnom/booking/pkg/config"
	"github.com/eighthGnom/booking/pkg/models"
)

var appConfig *config.AppConfig
var functions = template.FuncMap{}

// SetTemplatesConfig sets the config for the template package
func SetTemplatesConfig(config *config.AppConfig) {
	appConfig = config
}

// AddDefaultData adds to the templates default data
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate render templates and writes them to response
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
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

	td = AddDefaultData(td)

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
