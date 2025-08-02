package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID           uint64 `json:"user_id"`
	KratosIdentityID string `json:"kratos_identity_id"`
	Email            string `json:"email"`
	UserType         string `json:"user_type"`
	DisplayName      string `json:"display_name"`
	IsActive         bool   `json:"is_active"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey     []byte
	issuer        string
	tokenLifetime time.Duration
}

func NewJWTService(secretKey, issuer string, tokenLifetime time.Duration) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		issuer:        issuer,
		tokenLifetime: tokenLifetime,
	}
}

func (j *JWTService) GenerateToken(userID uint64, kratosIdentityID, email, userType, displayName string, isActive bool) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:           userID,
		KratosIdentityID: kratosIdentityID,
		Email:            email,
		UserType:         userType,
		DisplayName:      displayName,
		IsActive:         isActive,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   kratosIdentityID,
			Audience:  []string{"audora-api"},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.tokenLifetime)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTService) RefreshToken(claims *Claims) (string, error) {
	// Create new claims with updated expiration
	now := time.Now()
	newClaims := &Claims{
		UserID:           claims.UserID,
		KratosIdentityID: claims.KratosIdentityID,
		Email:            claims.Email,
		UserType:         claims.UserType,
		DisplayName:      claims.DisplayName,
		IsActive:         claims.IsActive,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   claims.Subject,
			Audience:  claims.Audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.tokenLifetime)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString(j.secretKey)
}