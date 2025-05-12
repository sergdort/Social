package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTAutheticator struct {
	secretKey string
	audience  string
	issuer    string
	tokenHost string
	expire    time.Duration
}

func NewJWTAutheticator(
	secretKey, audience, issuer, tokenHost string,
	expire time.Duration,
) *JWTAutheticator {
	return &JWTAutheticator{
		secretKey: secretKey,
		audience:  audience,
		issuer:    issuer,
		tokenHost: tokenHost,
		expire:    expire,
	}
}

func (auth *JWTAutheticator) GenerateToken(ctx context.Context, userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(auth.expire).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": auth.tokenHost,
		"aud": auth.tokenHost,
	}

	return auth.Generate(claims)
}

func (auth *JWTAutheticator) Generate(claims jwt.Claims) (string, error) {
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
