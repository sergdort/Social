package main

import (
	"github.com/sergdort/Social/internal/auth"
	"github.com/sergdort/Social/internal/store"
	"github.com/sergdort/Social/internal/store/cache"
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
		mockCacheStore := app.cache.Users.(*cache.MockUsersCache)
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

		mockCacheStore := app.cache.Users.(*cache.MockUsersCache)
		mockCacheStore.Mock.On("Get", mock.Anything, int64(userID)).Return(&store.User{ID: int64(userID)}, nil)

		mockUsersStore := app.store.Users.(*store.MockUsersRepository)
		mockUsersStore.Mock.On("GetByID", mock.Anything, int64(1)).Return(&store.User{ID: 1}, nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer fake_token")

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
