package main

import (
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/internal/auth"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t, config{
		redisCfg: redisConfig{
			enabled: true,
		},
	})

	mux := app.mount()

	t.Cleanup(func() {
		mockCacheStore := app.cache.Users.(*domain.MockUsersCache)
		mockCacheStore.Calls = nil
	})

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		fakeToken := "fake_token"
		userID := 42

		mockAuth := app.authenticator.(*auth.MockAuthenticator)
		mockAuth.Mock.On("ValidateToken", fakeToken).Return(makeFakeToken(fakeToken, userID), nil)

		mockCacheStore := app.cache.Users.(*domain.MockUsersCache)
		mockCacheStore.Mock.On("Get", mock.Anything, int64(userID)).Return(&domain.User{ID: int64(userID)}, nil)

		mockUsersStore := app.store.Users.(*domain.MockUsersRepository)
		mockUsersStore.Mock.On("GetByID", mock.Anything, int64(1)).Return(&domain.User{ID: 1}, nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer fake_token")

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
