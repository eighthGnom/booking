package handlers

import (
	"net/http"

	"github.com/eighthGnom/booking/pkg/config"
	"github.com/eighthGnom/booking/pkg/models"
	"github.com/eighthGnom/booking/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	AppConfig *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(appConfig *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: appConfig,
	}
}

// SetRepoForHandlers sets the repository for the handlers
func SetRepoForHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	rp.AppConfig.SessionManager.Put(r.Context(), "user_ip", r.RemoteAddr)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	userIP := rp.AppConfig.SessionManager.GetString(r.Context(), "user_ip")
	stringMap["test"] = "This is tested message for about page"
	stringMap["user_ip"] = userIP
	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringData: stringMap,
	})

}
