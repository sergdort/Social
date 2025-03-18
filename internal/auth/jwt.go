package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAutheticator struct {
	secretKey string
	audience  string
	issuer    string
}

func NewJWTAutheticator(secretKey, audience, issuer string) *JWTAutheticator {
	return &JWTAutheticator{
		secretKey: secretKey,
		audience:  audience,
		issuer:    issuer,
	}
}

func (auth *JWTAutheticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(auth.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (auth *JWTAutheticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(auth.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(auth.audience),
		jwt.WithIssuer(auth.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
