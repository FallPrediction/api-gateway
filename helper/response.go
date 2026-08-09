package helper

import (
	"api-gateway/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func writeJSON(w http.ResponseWriter, obj any) error {
	jsonBytes, err := json.Marshal(obj)
	logger := logger.NewLogger()
	if err != nil {
		logger.Error(
			"Marshal JSON failed.",
			zap.String("obj", fmt.Sprintf("%+v", obj)),
			zap.String("err", err.Error()),
		)
		return err
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		logger.Error(
			"Response writer writes data failed.",
			zap.String("json_bytes", string(jsonBytes)),
			zap.String("err", err.Error()),
		)
	}
	return err
}

func JSONResponse(w http.ResponseWriter, code int, headers map[string]string, obj any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(code)
	if err := writeJSON(w, obj); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
