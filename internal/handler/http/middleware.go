package handlerhttp

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/g6834/team21/authproto.git/pkg/tokensservice"
)

type Profile struct {
	UserId string
}
type ctxKey struct{}

func (c ctxKey) String() string {
	return "profile"
}

func (h *HandlerHttp) AuthValidate(ctx context.Context, w http.ResponseWriter, r *http.Request) (Profile, error) {

	tokenAccess, err := r.Cookie("access_token")
	if err != nil {
		return Profile{}, err
	}

	tokenRefresh, err := r.Cookie("refresh_token")
	if err != nil {
		return Profile{}, err
	}

	token, err := h.auth.Tokensservice.Validate(ctx, &tokensservice.Tokens{
		Access:  tokenAccess.Value,
		Refresh: tokenRefresh.Value,
	})
	if err != nil {
		return Profile{}, err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.Access,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.Refresh,
		HttpOnly: true,
	})

	h.logger.Info(token.UserID)

	return Profile{
		UserId: token.UserID,
	}, nil
}
func (h *HandlerHttp) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()
		profile, err := h.AuthValidate(ctx, w, r)
		if err != nil {
			h.logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return

		}
		r = r.WithContext(context.WithValue(ctx, ctxKey{}, profile))
		h.logger.Info(profile.UserId)
		next.ServeHTTP(w, r)
	})
}

func (h *HandlerHttp) CheckProfiler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.profiler.Status() {
			h.logger.Debug("Profiler access off")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
