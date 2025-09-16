package middleware

import (
	"api-gateway/helper"
	"api-gateway/logger"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var _ Middleware = new(Recovery)

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
				helper.JSONResponse(w, http.StatusInternalServerError, map[string]string{
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
