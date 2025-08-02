// internal/auth/application/service.go
package application

import (
	"fmt"
	"music-app-backend/internal/auth/adapters/repository"
	appError "music-app-backend/pkg/error"
	"music-app-backend/pkg/jwt"
	"music-app-backend/pkg/kratos"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepo     *repository.AuthRepository
	kratosClient *kratos.Client
	jwtService   *jwt.JWTService
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"`
	User        UserInfo `json:"user"`
}

type UserInfo struct {
	ID               uint64 `json:"id"`
	KratosIdentityID string `json:"kratos_identity_id"`
	Email            string `json:"email"`
	DisplayName      string `json:"display_name"`
	UserType         string `json:"user_type"`
	IsActive         bool   `json:"is_active"`
}

func NewAuthService(authRepo *repository.AuthRepository, kratosClient *kratos.Client, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		authRepo:     authRepo,
		kratosClient: kratosClient,
		jwtService:   jwtService,
	}
}

// VerifySessionAndIssueJWT validates Kratos session and issues Audora JWT
func (s *AuthService) VerifySessionAndIssueJWT(sessionToken string) (*LoginResponse, error) {
	// Step 1: Verify session with Kratos
	session, err := s.kratosClient.VerifySession(sessionToken)
	if err != nil {
		if kratosErr, ok := err.(*kratos.KratosError); ok {
			return nil, appError.NewUnauthorizedError(err, kratosErr.Message)
		}
		return nil, appError.NewInternalError(err, "failed to verify session with Kratos")
	}

	// Step 2: Find user in our database using Kratos identity ID
	kratosIdentityID, err := uuid.Parse(session.Identity.ID)
	if err != nil {
		return nil, appError.NewBadRequestError(err, "invalid kratos identity ID")
	}

	user, err := s.authRepo.FindUserByKratosIdentityID(kratosIdentityID)
	if err != nil {
		return nil, appError.NewNotFoundError(err, "user not found in Audora database")
	}

	// Step 3: Validate user is active
	if !user.IsActive {
		return nil, appError.NewForbiddenError(nil, "user account is deactivated")
	}

	// Step 4: Update last login time
	now := time.Now().Unix()
	user.LastLoginAt = &now
	if err := s.authRepo.UpdateUserLastLogin(user); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Step 5: Generate Audora JWT
	tokenLifetime := 24 * time.Hour // 24 hours
	accessToken, err := s.jwtService.GenerateToken(
		user.ID,
		user.KratosIdentityID.String(),
		user.Email,
		user.UserType,
		user.DisplayName,
		user.IsActive,
	)
	if err != nil {
		return nil, appError.NewInternalError(err, "failed to generate access token")
	}

	// Step 6: Return response
	return &LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(tokenLifetime.Seconds()),
		User: UserInfo{
			ID:               user.ID,
			KratosIdentityID: user.KratosIdentityID.String(),
			Email:            user.Email,
			DisplayName:      user.DisplayName,
			UserType:         user.UserType,
			IsActive:         user.IsActive,
		},
	}, nil
}

// ValidateJWT validates an Audora JWT token
func (s *AuthService) ValidateJWT(tokenString string) (*jwt.Claims, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, appError.NewUnauthorizedError(err, "invalid or expired token")
	}

	// Optional: Check if user is still active in database
	kratosIdentityID, err := uuid.Parse(claims.KratosIdentityID)
	if err != nil {
		return nil, appError.NewUnauthorizedError(err, "invalid identity ID in token")
	}

	user, err := s.authRepo.FindUserByKratosIdentityID(kratosIdentityID)
	if err != nil {
		return nil, appError.NewUnauthorizedError(err, "user not found")
	}

	if !user.IsActive {
		return nil, appError.NewForbiddenError(nil, "user account is deactivated")
	}

	return claims, nil
}

// RefreshToken refreshes an existing JWT token
func (s *AuthService) RefreshToken(tokenString string) (*LoginResponse, error) {
	// Validate existing token first
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Generate new token
	newToken, err := s.jwtService.RefreshToken(claims)
	if err != nil {
		return nil, appError.NewInternalError(err, "failed to refresh token")
	}

	tokenLifetime := 24 * time.Hour
	return &LoginResponse{
		AccessToken: newToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(tokenLifetime.Seconds()),
		User: UserInfo{
			ID:               claims.UserID,
			KratosIdentityID: claims.KratosIdentityID,
			Email:            claims.Email,
			DisplayName:      claims.DisplayName,
			UserType:         claims.UserType,
			IsActive:         claims.IsActive,
		},
	}, nil
}

// GetCurrentUser returns user info from JWT claims
func (s *AuthService) GetCurrentUser(claims *jwt.Claims) (*UserInfo, error) {
	return &UserInfo{
		ID:               claims.UserID,
		KratosIdentityID: claims.KratosIdentityID,
		Email:            claims.Email,
		DisplayName:      claims.DisplayName,
		UserType:         claims.UserType,
		IsActive:         claims.IsActive,
	}, nil
}
