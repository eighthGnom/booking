package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	newHandler := nosurf.New(next)

	newHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return newHandler
}

// SaveSession loads and saves the session on every request
func SaveSession(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}
