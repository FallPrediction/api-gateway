package handler

import (
	"net/http"
	"sync/atomic"

	"github.com/FallPrediction/api-gateway/internal/response"
)

var _ Handler = (*Gateway)(nil)

type Health struct {
	isShuttingDown *atomic.Bool
}

func (h *Health) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.isShuttingDown.Load() {
			response.JSONResponse(w, http.StatusServiceUnavailable, map[string]string{}, map[string]string{
				"msg": "Shutting down",
			})
			return
		}
		response.JSONResponse(w, http.StatusOK, map[string]string{}, map[string]string{
			"msg": "OK",
		})
	})
}

func NewHealth(isShuttingDown *atomic.Bool) Health {
	return Health{isShuttingDown}
}
