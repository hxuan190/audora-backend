package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ctx2 "music-app-backend/pkg/context"
	analyticsModule "music-app-backend/internal/analytics"
	musicModule "music-app-backend/internal/music"
	playbackModule "music-app-backend/internal/playback"
	socialModule "music-app-backend/internal/social"
	userModule "music-app-backend/internal/user"
	"music-app-backend/pkg/database"

	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	// In Docker containers, environment variables are passed directly
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables from container: %v", err)
	}

	IDconfig := goflakeid.NewConfig(1, 1, 0).WithAutoMachineID()
	generator, err := goflakeid.NewGenerator(*IDconfig)
	if err != nil {
		log.Fatalf("Failed to initialize ID generator: %v", err)
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
	api := router.Group("api")
	v1 := api.Group("v1")

	// Init Service Context
	serviceContext := ctx2.NewerviceContext(db.GetDB(), router, generator)

	// Module registration
	musicModule := musicModule.NewMusicModule(db.GetDB(), serviceContext)
	musicModule.RegisterRoutes(v1)

	userModule := userModule.NewUserModule(serviceContext, musicModule.Service)
	userModule.RegisterRoutes(v1)

	analyticsModule := analyticsModule.NewAnalyticsModule(db.GetDB())
	analyticsModule.RegisterRoutes(v1)

	playbackModule := playbackModule.NewPlaybackModule(db.GetDB())
	playbackModule.RegisterRoutes(v1)

	socialModule := socialModule.NewSocialModule(db.GetDB())
	socialModule.RegisterRoutes(v1)

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
