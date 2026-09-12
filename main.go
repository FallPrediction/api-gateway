package main

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/FallPrediction/api-gateway/internal/circuitbreaker"
	"github.com/FallPrediction/api-gateway/internal/handler"
	"github.com/FallPrediction/api-gateway/internal/logger"
	"github.com/FallPrediction/api-gateway/internal/metric"
	"github.com/FallPrediction/api-gateway/internal/middleware"
	"github.com/FallPrediction/api-gateway/internal/ratelimit"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var isShuttingDown atomic.Bool

const (
	shutdownPeriod      = 30 * time.Second
	shutdownHardPeriod  = 3 * time.Second
	readinessDrainDelay = 5 * time.Second
)

func getHandler(handler handler.Handler, middlewares ...middleware.Middleware) http.Handler {
	if len(middlewares) == 0 {
		return handler.Handle()
	}
	for i := 0; i < len(middlewares)-1; i++ {
		middlewares[i].SetNext(middlewares[i+1].Handle())
	}
	middlewares[len(middlewares)-1].SetNext(handler.Handle())
	return middlewares[0].Handle()
}

func main() {
	logger := logger.NewLogger()

	// SIGINT/SIGTERM context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	setRoute(mux, logger)

	// Set global context to server.
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr: ":8080",
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic("Fail to listen and serve", zap.String("err", err.Error()))
		}
	}()

	// 等待 SIGTERM/SIGINT
	<-rootCtx.Done()
	stop()

	shutdown(server, stopOngoingGracefully, logger)
}

func setRoute(mux *http.ServeMux, logger *zap.Logger) {
	mux.Handle("/", gatewayHandler(logger))
	mux.Handle("/metrics", promhttp.HandlerFor(metric.NewRegistry(), promhttp.HandlerOpts{}))
	healthHandler := handler.NewHealth(&isShuttingDown)
	mux.Handle("/healthz", healthHandler.Handle())
}

func gatewayHandler(logger *zap.Logger) http.Handler {
	proxy := handler.NewGateway()
	recoveryMiddleware := middleware.NewRecover()
	upstreams, err := upstream.LoadConfig("config.yaml")
	if err != nil {
		logger.Panic("Fail load config", zap.String("err", err.Error()))
	}
	upstreamMiddleware := middleware.NewUpstream(upstreams)
	corsMiddleware := middleware.NewCors()
	rateLimitMiddleware := middleware.NewRateLimit(ratelimit.NewRateLimiters(upstreams))
	circuitBreakerMiddleware := middleware.NewCircuitBreaker(circuitbreaker.NewNewCircuitBreakers(upstreams))
	logMiddleware := middleware.NewLog()
	authenticateMiddleware := middleware.NewAuthenticate()
	return getHandler(
		&proxy,
		&recoveryMiddleware,
		&upstreamMiddleware,
		&corsMiddleware,
		&rateLimitMiddleware,
		&circuitBreakerMiddleware,
		&logMiddleware,
		&authenticateMiddleware,
	)
}

func shutdown(server *http.Server, stopOngoingGracefully context.CancelFunc, logger *zap.Logger) {
	// isShuttingDown 設為 true，讓 health check API 回傳 503
	isShuttingDown.Store(true)
	logger.Info("Received shutdown signal, shutting down.")
	// 等待 5 秒讓 ALB/K8S 通過 health check API 感知到服務正在 shutdown
	time.Sleep(readinessDrainDelay)
	logger.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

	// 設定留 30 秒處理剩餘的請求
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()
	err := server.Shutdown(shutdownCtx)
	// 通知所有 handler global context 已取消
	stopOngoingGracefully()
	if err != nil {
		logger.Info("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		// 留最後 3 秒後強制結束
		time.Sleep(shutdownHardPeriod)
	}
	logger.Info("Server shut down gracefully.")
}
