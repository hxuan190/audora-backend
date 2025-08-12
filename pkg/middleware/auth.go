// pkg/middleware/auth.go
package middleware

import (
	"music-app-backend/internal/auth/application"
	jsonResponse "music-app-backend/pkg/json"
	"music-app-backend/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *application.AuthService
}

func NewAuthMiddleware(authService *application.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware that validates JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			jsonResponse.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			jsonResponse.ResponseBadRequest(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Validate token
		claims, err := m.authService.ValidateJWT(tokenString)
		if err != nil {
			jsonResponse.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("user_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("kratos_identity_id", claims.KratosIdentityID)
		c.Set("user_type", claims.UserType)
		c.Set("user_email", claims.Email)
		c.Set("user_tier", "free") // Default to free tier, can be updated later

		c.Next()
	}
}

// RequireArtist middleware that requires user to be an artist
func (m *AuthMiddleware) RequireArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user_claims")
		if !exists {
			jsonResponse.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		userClaims := claims.(*jwt.Claims)
		if userClaims.UserType != "artist" {
			jsonResponse.ResponseForbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireListener middleware that requires user to be a listener
func (m *AuthMiddleware) RequireListener() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user_claims")
		if !exists {
			jsonResponse.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		userClaims := claims.(*jwt.Claims)
		if userClaims.UserType != "listener" {
			jsonResponse.ResponseForbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware that optionally validates JWT token
// Useful for endpoints that work for both authenticated and anonymous users
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		// Try to validate token, but don't fail if invalid
		claims, err := m.authService.ValidateJWT(tokenString)
		if err == nil {
			c.Set("user_claims", claims)
			c.Set("user_id", claims.UserID)
			c.Set("kratos_identity_id", claims.KratosIdentityID)
			c.Set("user_type", claims.UserType)
			c.Set("user_email", claims.Email)
		}

		c.Next()
	}
}