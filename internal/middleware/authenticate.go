package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/FallPrediction/api-gateway/internal/auth"
	"github.com/FallPrediction/api-gateway/internal/response"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/golang-jwt/jwt/v5"
)

var _ Middleware = (*Authenticate)(nil)

type Authenticate struct {
	baseMiddleware
}

func (m *Authenticate) getUserClaim(r *http.Request) (auth.UserClaims, error) {
	token := r.Header.Get("authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return auth.UserClaims{}, errors.New("token missing")
	}
	claims, err := (&auth.Jwt{}).ParseUserToken(token[7:])
	if err != nil {
		return auth.UserClaims{}, err
	}
	return *claims, nil
}

func (m *Authenticate) setClaimToHeader(claim auth.UserClaims, r *http.Request) {
	r.Header.Set("UserId", strconv.Itoa(int(claim.UserId)))
	r.Header.Set("Email", claim.Email)
	r.Header.Set("Scope", claim.Scope)
}

func (m *Authenticate) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := upstream.GetUpstream(r.Context())
		if u != nil && u.Auth {
			userClaim, err := m.getUserClaim(r)
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.JSONResponse(w, http.StatusForbidden, map[string]string{}, map[string]string{
					"msg": "Token expired",
				})
				return
			} else if err != nil {
				response.JSONResponse(w, http.StatusUnauthorized, map[string]string{}, map[string]string{
					"msg": "Unauthorized",
				})
				return
			}
			if !slices.Contains(strings.Split(userClaim.Scope, ","), u.Name) {
				response.JSONResponse(w, http.StatusForbidden, map[string]string{}, map[string]string{
					"msg": "Forbidden",
				})
				return
			}
			m.setClaimToHeader(userClaim, r)
		}
		m.next.ServeHTTP(w, r)
	})
}

func NewAuthenticate() Authenticate {
	return Authenticate{}
}
