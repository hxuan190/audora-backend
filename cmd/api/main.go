package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	userModule "music-app-backend/internal/user"
	"music-app-backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	// In Docker containers, environment variables are passed directly
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables from container: %v", err)
	}

	user := os.Getenv("AUDORA_DB_USER")
	password := os.Getenv("AUDORA_DB_PASSWORD")
	dbname := os.Getenv("AUDORA_DB_NAME")
	host := os.Getenv("AUDORA_DB_HOST")
	port := os.Getenv("AUDORA_DB_PORT")

	// Use default values if not set
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}

	config := database.DefaultConfig()
	config.Host = host
	config.Port = port
	config.User = user
	config.Password = password
	config.DBName = dbname

	db, err := database.NewWithConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := gin.Default()

	userModule := userModule.NewUserModule(db.GetDB())
	userModule.RegisterRoutes(router)

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
