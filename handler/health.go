package handler

import "net/http"

var _ Handler = new(Health)

type Health struct {
	baseHandler
}

func (h *Health) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.JSONResponse(w, http.StatusOK, map[string]string{
			"msg": "OK",
		})
	})
}

func NewHealth() *Health {
	return &Health{}
}
