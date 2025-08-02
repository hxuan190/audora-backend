// internal/auth/adapters/http/handler.go
package http

import (
	"music-app-backend/internal/auth/application"
	appError "music-app-backend/pkg/error"
	jsonResponse "music-app-backend/pkg/json"
	"music-app-backend/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *application.AuthService
}

type LoginRequest struct {
	SessionToken string `json:"session_token" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if appErr, ok := appError.GetAppError(err); ok {
		jsonResponse.ResponseJSON(c, appErr.StatusCode, appErr.Message, appErr.Data)
		return true
	}

	jsonResponse.ResponseInternalError(c, err)
	return true
}

// Login validates Kratos session and issues Audora JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonResponse.ResponseBadRequest(c, "Invalid request: "+err.Error())
		return
	}

	response, err := h.authService.VerifySessionAndIssueJWT(request.SessionToken)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	jsonResponse.ResponseOK(c, response)
}

// LoginWithCookie validates Kratos session from cookie and issues Audora JWT
func (h *AuthHandler) LoginWithCookie(c *gin.Context) {
	// Extract session token from cookie
	sessionToken, err := c.Cookie("ory_kratos_session")
	if err != nil {
		jsonResponse.ResponseUnauthorized(c)
		return
	}

	response, err := h.authService.VerifySessionAndIssueJWT(sessionToken)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	jsonResponse.ResponseOK(c, response)
}

// RefreshToken refreshes an existing JWT token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonResponse.ResponseBadRequest(c, "Invalid request: "+err.Error())
		return
	}

	response, err := h.authService.RefreshToken(request.RefreshToken)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	jsonResponse.ResponseOK(c, response)
}

// Me returns current user information from JWT
func (h *AuthHandler) Me(c *gin.Context) {
	// Get claims from middleware context
	claims, exists := c.Get("user_claims")
	if !exists {
		jsonResponse.ResponseUnauthorized(c)
		return
	}

	userInfo, err := h.authService.GetCurrentUser(claims.(*jwt.Claims))
	if err != nil {
		h.HandleError(c, err)
		return
	}

	jsonResponse.ResponseOK(c, userInfo)
}

// Logout invalidates the current session (optional - mainly handled by frontend)
func (h *AuthHandler) Logout(c *gin.Context) {
	// Since we're using stateless JWT, logout is mainly handled client-side
	// But we can add token blacklisting here if needed in the future
	jsonResponse.ResponseOK(c, gin.H{"message": "Successfully logged out"})
}

// ValidateToken endpoint for other services to validate tokens
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		jsonResponse.ResponseUnauthorized(c)
		return
	}

	// Extract Bearer token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		jsonResponse.ResponseBadRequest(c, "Invalid authorization header format")
		return
	}

	claims, err := h.authService.ValidateJWT(tokenString)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	userInfo := &application.UserInfo{
		ID:               claims.UserID,
		KratosIdentityID: claims.KratosIdentityID,
		Email:            claims.Email,
		DisplayName:      claims.DisplayName,
		UserType:         claims.UserType,
		IsActive:         claims.IsActive,
	}

	jsonResponse.ResponseOK(c, userInfo)
}
