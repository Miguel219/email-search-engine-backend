package server

import (
	"fmt"
	"net/http"

	handlers "email-search-engine-backend/internal/server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func CreateServer() {
	router := chi.NewRouter()

	configServer(router)

	configEndpoints(router)

	fmt.Println("The server is running")
	http.ListenAndServe(":8000", router)
}

func configServer(router chi.Router) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func configEndpoints(router chi.Router) {
	//Emails
	router.Mount("/emails", configEmailsEndpoints())
}

// REST routes for "emails" resource
func configEmailsEndpoints() chi.Router {
	router := chi.NewRouter()
	router.Post("/list", handlers.ListEmails)
	router.Post("/search", handlers.SearchEmails) // GET /emails/search

	// r.Post("/", CreateEmail)       // POST /emails
	// r.Route("/{emailID}", func(r chi.Router) {
	// 	router.Use(EmailCtx)            // Load the *Email on the request context
	// 	router.Get("/", GetEmail)       // GET /emails/123
	// 	router.Put("/", UpdateEmail)    // PUT /emails/123
	// 	router.Delete("/", DeleteEmail) // DELETE /emails/123
	// })
	return router
}
