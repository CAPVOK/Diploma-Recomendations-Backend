package service

import (
	"time"

	"diprec_api/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type AuthService struct {
	config *JWTConfig
}

func NewAuthService(cfg *JWTConfig) *AuthService {
	return &AuthService{config: cfg}
}

func (a *AuthService) GenerateTokens(user *domain.User) (*domain.TokenPair, error) {
	// Access token
	accessClaims := jwt.MapClaims{
		"userID":    user.ID,
		"role":      user.Role,
		"tokenType": "access",
		"exp":       time.Now().Add(a.config.AccessExpiry).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(a.config.SecretKey))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"userID":    user.ID,
		"role":      user.Role,
		"tokenType": "refresh",
		"exp":       time.Now().Add(a.config.RefreshExpiry).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(a.config.SecretKey))
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresAt:    time.Now().Add(a.config.AccessExpiry),
	}, nil
}

func (a *AuthService) ValidateToken(tokenString string) (uint, string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, "", "", domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", "", domain.ErrUnauthorized
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, "", "", domain.ErrUnauthorized
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", "", domain.ErrUnauthorized
	}

	tokenType, _ := claims["tokenType"].(string)
	return uint(userID), role, tokenType, nil
}

func (a *AuthService) ValidateRefreshToken(tokenString string) (uint, string, error) {
	userID, role, tokenType, err := a.ValidateToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	if tokenType != "refresh" {
		return 0, "", domain.ErrInvalidTokenType
	}

	return userID, role, nil
}
