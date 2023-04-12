package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eighthGnom/booking/internal/config"
	"github.com/eighthGnom/booking/internal/driver"
	"github.com/eighthGnom/booking/internal/forms"
	"github.com/eighthGnom/booking/internal/helpers"
	"github.com/eighthGnom/booking/internal/models"
	"github.com/eighthGnom/booking/internal/render"
	"github.com/eighthGnom/booking/internal/repository"
	"github.com/eighthGnom/booking/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	AppConfig *config.AppConfig
	DB        repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(appConfig *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		AppConfig: appConfig,
		DB:        dbrepo.NewPostgresRepo(db.SQL, appConfig),
	}
}

// SetRepoForHandlers sets the repository for the handlers
func SetRepoForHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.AppConfig.SessionManager.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handle the posting the reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	form := forms.New(r.Form)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	m.AppConfig.SessionManager.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("The start of reservation is %s , the end of reservation is %s", start, end)))
}

type jsonMessage struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON sent json
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	message := jsonMessage{
		Ok:      true,
		Message: "Available",
	}
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.AppConfig.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.AppConfig.ErrorLog.Println("Can't get value form session")
		m.AppConfig.SessionManager.Put(r.Context(), "error", "Cannot get item from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.AppConfig.SessionManager.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})

}
