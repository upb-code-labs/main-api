package infrastructure

import (
	"errors"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/session/domain"
	"github.com/golang-jwt/jwt/v4"
)

type JwtTokenHandler struct {
	Secret          string
	ExpirationHours int
}

var jwtTokenHandler *JwtTokenHandler

func GetJwtTokenHandler() *JwtTokenHandler {
	if jwtTokenHandler == nil {
		jwtTokenHandler = &JwtTokenHandler{
			Secret:          GetEnvironment().JwtSecret,
			ExpirationHours: GetEnvironment().JwtExpirationHours,
		}
	}

	return jwtTokenHandler
}

type Claims struct {
	UUID string `json:"uuid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (handler JwtTokenHandler) GenerateToken(user entities.User) (string, error) {
	claims := Claims{
		UUID: user.UUID,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(handler.ExpirationHours) * time.Hour),
			),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "codelabs",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(handler.Secret))
	return signedToken, err
}

func (handler JwtTokenHandler) ValidateToken(token string) (domain.JwtCustomClaims, error) {
	claims := Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.Secret), nil
	})
	if err != nil {
		return domain.JwtCustomClaims{}, err
	}

	if !parsedToken.Valid {
		return domain.JwtCustomClaims{}, errors.New("invalid token")
	}

	return domain.JwtCustomClaims{
		UUID: claims.UUID,
		Role: claims.Role,
	}, nil
}
