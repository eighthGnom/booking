package models

import "github.com/eighthGnom/booking/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringData map[string]string
	IntData    map[string]int
	FloatData  map[string]float32
	Data       map[string]interface{}
	CSRFToken  string
	Flash      string
	Warning    string
	Error      string
	FormValues *forms.Form
}
