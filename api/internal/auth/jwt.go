package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey            string
	tokenLifetimeMinutes int
}

func NewJWTService(secretKey string, tokenLifetimeMinutes int) *JWTService {
	return &JWTService{
		secretKey:            secretKey,
		tokenLifetimeMinutes: tokenLifetimeMinutes,
	}
}

type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(s.tokenLifetimeMinutes) * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(s.secretKey), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
