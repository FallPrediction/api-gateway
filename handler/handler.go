package handler

import (
	"api-gateway/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Handler interface {
	Handle() http.Handler
}

type baseHandler struct{}

func (h *baseHandler) writeJSON(w http.ResponseWriter, obj any) error {
	jsonBytes, err := json.Marshal(obj)
	logger := logger.NewLogger()
	if err != nil {
		logger.Info(
			"Marshal JSON failed.",
			zap.String("obj", fmt.Sprintf("%+v", obj)),
			zap.String("err", err.Error()),
		)
		return err
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		logger.Info(
			"Response writer writes data failed.",
			zap.String("json_bytes", string(jsonBytes)),
			zap.String("err", err.Error()),
		)
	}
	return err
}

func (h *baseHandler) JSONResponse(w http.ResponseWriter, code int, obj any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := h.writeJSON(w, obj); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
