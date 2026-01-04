package main

import (
	"fmt"
	"log"

	"bookmgr/config"
	"bookmgr/database"
	"bookmgr/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Book API Server...")
	log.Printf("Database: %s@%s:%d/%s", cfg.DatabaseUser, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
	log.Printf("Skip Bootstrap: %v", cfg.SkipBootstrap)

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	handler := handlers.New(db)

	// Setup Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		api.POST("/books", handler.CreateBook)
		api.GET("/books", handler.GetAllBooks)
		api.GET("/books/:id", handler.GetBook)
		api.PUT("/books/:id", handler.UpdateBook)
		api.DELETE("/books/:id", handler.DeleteBook)
	}

	// Serve static files (Web UI)
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on port %d", cfg.ServerPort)
	log.Printf("Web UI: http://localhost:%d", cfg.ServerPort)
	log.Printf("API: http://localhost:%d/api/books", cfg.ServerPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
