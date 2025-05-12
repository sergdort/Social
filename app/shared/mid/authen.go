package mid

import (
	"context"
	"github.com/sergdort/Social/app/shared/errs"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
	"strings"
)

// Bearer processes JWT authentication logic.
func Bearer(ath *domain.AuthUseCase) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return errs.Newf(errs.Unauthenticated, "expected authorization header format: Bearer <token>")
			}

			// "Bearer <token>"
			token := authHeader[7:]
			calaims, err := ath.ValidateToken(ctx, token)

			if err != nil {
				return errs.Newf(errs.Unauthenticated, "invalid token")
			}

			ctx = setUserID(ctx, calaims.UserID)

			return next(ctx, r)
		}

		return h
	}

	return m
}
