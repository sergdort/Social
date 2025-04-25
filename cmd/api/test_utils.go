package main

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/store"
	"github.com/sergdort/Social/business/platform/store/cache"
	"github.com/sergdort/Social/internal/auth"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestApplication(t *testing.T, cfg config) *application {
	return &application{
		config: cfg,
		store: store.Storage{
			Posts:    domain.NewMockPostsRepository(t),
			Users:    domain.NewMockUsersRepository(t),
			Comments: domain.NewMockCommentsRepository(t),
			Follows:  domain.NewMockFollowsRepository(t),
			Roles:    domain.NewMockRolesRepository(t),
		},
		logger:        zap.NewNop().Sugar(),
		mailer:        nil,
		authenticator: auth.NewMockAuthenticator(t),
		cache: cache.Storage{
			Users: domain.NewMockUsersCache(t),
		},
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

func makeFakeToken(token string, userID int) *jwt.Token {
	// Create some dummy claims
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"role": "admin",
	}

	// Create a fake token object (not signed or verified)
	return &jwt.Token{
		Raw:    token,
		Claims: claims,
		Method: jwt.SigningMethodHS256,
		Valid:  true, // Set manually if you're skipping Parse/Verify
	}
}
