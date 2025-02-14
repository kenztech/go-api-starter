package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kenztech/go-api-starter/handlers"
	"github.com/kenztech/go-api-starter/models"
	"github.com/kenztech/go-api-starter/utils"
)

func main() {
	db, err := models.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return
	}
	defer func() {
		if err := db.Client().Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
	}()
	log.Println("Database connection established:", db.Name())

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link", "X-Total-Count", "Set-Cookie"},
		AllowCredentials: true,
	}))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.InitRoutes(r, db)

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Server starting at port %v", port)
	http.ListenAndServe(":"+port, r)
}
