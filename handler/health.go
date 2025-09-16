package handler

import (
	"api-gateway/helper"
	"net/http"
)

var _ Handler = new(Health)

type Health struct{}

func (h *Health) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helper.JSONResponse(w, http.StatusOK, map[string]string{}, map[string]string{
			"msg": "OK",
		})
	})
}

func NewHealth() *Health {
	return &Health{}
}
