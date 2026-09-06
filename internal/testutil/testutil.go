package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/FallPrediction/api-gateway/internal/auth"
	"github.com/FallPrediction/api-gateway/internal/upstream"
	"github.com/golang-jwt/jwt/v5"
)

func CreateRequestWithUpstream(u *upstream.Upstream) *http.Request {
	req := httptest.NewRequest("GET", "/"+u.Name, nil)
	req = req.WithContext(upstream.WithUpstream(req.Context(), u))
	return req
}

func MockResponse(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})
}

func CreateToken(scope string, expiresAt time.Time) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = &auth.UserClaims{
		UserData: auth.UserData{
			UserId: 1,
			Email:  "johndoe@example.com",
			Scope:  scope,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	return t.SignedString([]byte(os.Getenv("APP_KEY")))
}
