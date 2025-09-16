package middleware

import (
	"api-gateway/helper"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var _ Middleware = new(Authenticate)

type Authenticate struct {
	baseMiddleware
}

func (m *Authenticate) getUserClaim(r *http.Request) (helper.UserClaims, bool) {
	token := r.Header.Get("authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return helper.UserClaims{}, false
	}
	claims, err := (&helper.Jwt{}).ParseUserToken(token[7:])
	if err != nil || claims.ExpiresAt == nil {
		return helper.UserClaims{}, false
	}
	return *claims, true
}

func (m *Authenticate) setClaimToHeader(claim helper.UserClaims, r *http.Request) {
	r.Header.Set("UserId", strconv.Itoa(int(claim.UserId)))
	r.Header.Set("Email", claim.Email)
	r.Header.Set("Scope", claim.Scope)
}

func (m *Authenticate) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, ok := helper.GetOrigin(r.URL.Path)
		if ok && route.Auth {
			userClaim, ok := m.getUserClaim(r)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if userClaim.ExpiresAt.Compare(time.Now()) == -1 {
				http.Error(w, "Token expired", http.StatusForbidden)
				return
			}
			if !slices.Contains(strings.Split(userClaim.Scope, ","), route.Name) {
				http.Error(w, "Forbidden", http.StatusForbidden)
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
