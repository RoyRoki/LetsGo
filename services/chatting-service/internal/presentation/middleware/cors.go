package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

// CORS returns a middleware handler for Cross-Origin Resource Sharing
func CORS() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"}, // TODO: Replace with real domain in prod
		AllowedMethods:   []string{"GET", "PUT", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	}).Handler
}
