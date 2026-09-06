package middleware

import (
	"fmt"
	"net/http"

	"github.com/FallPrediction/api-gateway/internal/logger"
	"github.com/FallPrediction/api-gateway/internal/response"

	"go.uber.org/zap"
)

var _ Middleware = (*Recovery)(nil)

type Recovery struct {
	baseMiddleware
}

func (m *Recovery) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := logger.NewLogger()
				logger.Error(
					"Panic recovered.",
					zap.String("err", fmt.Sprintf("%v", err)),
				)
				response.JSONResponse(w, http.StatusInternalServerError, map[string]string{}, map[string]string{
					"msg": "Internal Server Error",
				})
			}
		}()

		m.next.ServeHTTP(w, r)
	})
}

func NewRecover() Recovery {
	return Recovery{}
}
