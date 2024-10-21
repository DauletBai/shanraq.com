package routes

import (
	"github.com/go-chi/chi/v5"
	"shanraq.com/internal/handlers"
	"shanraq.com/internal/middleware"
)

func SetupRoutes(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)

	r.Group(func(r chi.Router)  {
		r.Use(middleware.Authenticate)

		r.Route("/articles", func(r chi.Router)  {
			r.Get("/", handlers.GetArticles)

			r.Group(func(r chi.Router)  {
				r.Use(middleware.Authorize("admin", "editor", "author"))
				r.Post("/", handlers.CreateArticle)
				r.Put("{id}", handlers.UpdateArticle)
				r.Delete("/{id}", handlers.DeleteArticle)
			})
		})

		// Other routers
	})

	return r
}