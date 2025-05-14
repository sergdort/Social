package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sergdort/Social/business/domain"
	"strconv"
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
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(auth.expire).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": auth.tokenHost,
		"aud": auth.tokenHost,
	}

	return auth.generate(claims)
}

func (auth *JWTAutheticator) ValidateToken(ctx context.Context, token string) (domain.Claims, error) {
	jwtToken, err := auth.validate(token)
	if err != nil {
		return domain.Claims{}, err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		return domain.Claims{}, err
	}
	return domain.Claims{userID}, nil
}

func (auth *JWTAutheticator) generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(auth.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (auth *JWTAutheticator) validate(token string) (*jwt.Token, error) {
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
