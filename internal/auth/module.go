// internal/auth/module.go
package auth

import (
	"music-app-backend/internal/auth/adapters/http"
	"music-app-backend/internal/auth/adapters/repository"
	"music-app-backend/internal/auth/application"
	"music-app-backend/pkg/jwt"
	"music-app-backend/pkg/kratos"
	"music-app-backend/pkg/middleware"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthModule struct {
	Repository   *repository.AuthRepository
	Service      *application.AuthService
	Handler      *http.AuthHandler
	Middleware   *middleware.AuthMiddleware
	KratosClient *kratos.Client
	JWTService   *jwt.JWTService
}

func NewAuthModule(db *gorm.DB) *AuthModule {
	// Initialize Kratos client
	kratosPublicURL := os.Getenv("KRATOS_PUBLIC_URL")
	kratosAdminURL := os.Getenv("KRATOS_ADMIN_URL")
	if kratosPublicURL == "" {
		kratosPublicURL = "http://localhost:4433"
	}
	if kratosAdminURL == "" {
		kratosAdminURL = "http://localhost:4434"
	}
	kratosClient := kratos.NewClient(kratosPublicURL, kratosAdminURL)

	// Initialize JWT service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-jwt-secret-here" // Default for development
	}
	jwtService := jwt.NewJWTService(jwtSecret, "audora-api", 24*time.Hour)

	// Initialize repository, service, handler, and middleware
	authRepo := repository.NewAuthRepository(db)
	authService := application.NewAuthService(authRepo, kratosClient, jwtService)
	authHandler := http.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	return &AuthModule{
		Repository:   authRepo,
		Service:      authService,
		Handler:      authHandler,
		Middleware:   authMiddleware,
		KratosClient: kratosClient,
		JWTService:   jwtService,
	}
}

func (a *AuthModule) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("auth")
	{
		// Public endpoints
		auth.POST("/login", a.Handler.Login)                  // Login with session token from body
		auth.POST("/login/cookie", a.Handler.LoginWithCookie) // Login with session token from cookie
		auth.POST("/refresh", a.Handler.RefreshToken)         // Refresh JWT token
		auth.POST("/validate", a.Handler.ValidateToken)       // Validate token (for other services)

		// Protected endpoints
		protected := auth.Group("")
		protected.Use(a.Middleware.RequireAuth())
		{
			protected.GET("/me", a.Handler.Me)          // Get current user info
			protected.POST("/logout", a.Handler.Logout) // Logout (optional)
		}
	}
}
