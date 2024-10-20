package main

import (
	"log"
	"net/http"

	"shanraq.com/internal/config"
	"shanraq.com/internal/database"
	"shanraq.com/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error load config: %v", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Database connection errors: %v", err)
	}
	defer db.Close()

	router := routes.SetupRoutes(db)

	log.Printf("The server is running on the port %s", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, router); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}