package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
	Values url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(errors),
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	v := r.Form.Get(field)
	if v == "" {
		f.Errors.Add(field, "field cannot be blank")
		return false
	}
	return true
}
