package auth

import "github.com/golang-jwt/jwt/v5"

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
	return nil, nil
}
