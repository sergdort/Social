package mid

import (
	"context"
	"encoding/base64"
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

			ctx = setAuthUserID(ctx, calaims.UserID)

			return next(ctx, r)
		}

		return h
	}

	return m
}

func Basic(username string, pass string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Basic ") {
				return errs.Newf(errs.Unauthenticated, "missing authorization header")
			}
			// "Basic tokenstringhere"
			base64Token := authHeader[6:]
			decodedToken, err := base64.StdEncoding.DecodeString(base64Token)
			if err != nil {
				return errs.Newf(errs.Unauthenticated, "invalid token %s", err.Error())
			}

			creds := strings.SplitN(string(decodedToken), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				return errs.Newf(errs.Unauthenticated, "invalid token")
			}

			return next(ctx, r)
		}
		return h
	}
	return m
}
