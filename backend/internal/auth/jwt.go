package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"taladsync/internal/domain"
)

const tokenTTL = 24 * time.Hour

type Claims struct {
	UserID uint        `json:"userId"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenService interface {
	Sign(userID uint, role domain.Role) (string, error)
	Parse(token string) (*Claims, error)
}

type tokenService struct {
	secret []byte
}

func NewTokenService(secret []byte) TokenService {
	return &tokenService{secret: secret}
}

func (s *tokenService) Sign(userID uint, role domain.Role) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}


func (s *tokenService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
