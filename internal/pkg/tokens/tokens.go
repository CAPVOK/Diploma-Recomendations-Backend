package tokens

import (
	"diprec_api/internal/domain"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTConfig struct {
	SecretKey        string `json:"secretKey"`
	AccessExpiresAt  int64  `json:"accessExpiresAt"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
}

func GenerateTokens(user domain.User) (*domain.TokenPair, error) {
	accessClaims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte("super_secret_key"))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte("super_secret_key"))
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresAt:    time.Now().Add(time.Hour * 24),
	}, nil
}
