package middleware

import (
	"io"
	"math/rand/v2"
	"net/http"

	"github.com/FallPrediction/api-gateway/internal/logger"
	"github.com/FallPrediction/api-gateway/internal/response"

	"go.uber.org/zap"
)

var _ Middleware = (*Log)(nil)

type Log struct {
	baseMiddleware
}

func (m *Log) getTraceId() string {
	resultBytes := [5]byte{}
	character := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range len(resultBytes) {
		resultBytes[i] = character[rand.IntN(len(character))]
	}
	return string(resultBytes[:])
}

func (m *Log) Handle() http.Handler {
	l := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := response.NewResponseWriter(w)
		body, _ := io.ReadAll(r.Body)
		traceId := m.getTraceId()
		l.Info(
			"",
			zap.String("trace_id", traceId),
			zap.String("type", "request"),
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path),
			zap.String("host", r.Host),
			zap.Object("header", (*logger.ZapHeader)(&r.Header)),
			zap.String("body", string(body)),
		)
		m.next.ServeHTTP(rw, r)
		l.Info(
			"",
			zap.String("trace_id", traceId),
			zap.String("type", "response"),
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path),
			zap.String("host", r.Host),
			zap.Object("header", (*logger.ZapHeader)(&r.Header)),
			zap.String("body", string(rw.Body)),
			zap.Int("status_code", rw.StatusCode),
		)
	})
}

func NewLog() Log {
	return Log{}
}
