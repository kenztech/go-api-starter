package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kenztech/go-api-starter/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitRoutes(r *chi.Mux, db *mongo.Database) {
	authHandler := NewAuthHandler(db)
	userHandler := NewUserHandler(db)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)
			r.With(middlewares.Authenticate).Get("/me", authHandler.Me)
			r.With(middlewares.Authenticate).Post("/logout", authHandler.Logout)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middlewares.Authenticate)
			r.Use(middlewares.AdminOnly)

			r.Get("/", userHandler.GetUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Post("/", userHandler.CreateUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})
}
