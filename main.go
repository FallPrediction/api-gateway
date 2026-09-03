package main

import (
	"api-gateway/handler"
	"api-gateway/helper"
	"api-gateway/logger"
	"api-gateway/middleware"
	"context"
	"net"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var isShuttinDown atomic.Bool

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

	upstreams, err := helper.LoadConfig("config.yaml")
	if err != nil {
		logger.Panic("Fail load config", zap.String("err", err.Error()))
	}
	helper.Upstreams = upstreams

	// SIGINT/SIGTERM context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	setRoute()

	// Set global context to server.
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr: ":8080",
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// 等待 SIGTERM/SIGINT
	<-rootCtx.Done()
	stop()

	// isShuttinDown 設為 true，讓 health check API 回傳 503
	isShuttinDown.Store(true)
	logger.Info("Received shutdown signal, shutting down.")
	// 等待 5 秒讓 ALB/K8S 通過 health check API 感知到服務正在 shutdown
	time.Sleep(readinessDrainDelay)
	logger.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

	// 設定留 30 秒處理剩餘的請求
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	// 通知所有 handler global context 已取消
	stopOngoingGracefully()
	if err != nil {
		logger.Info("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		// 留最後 3 秒後強制結束
		time.Sleep(shutdownHardPeriod)
	}
	logger.Info("Server shut down gracefully.")
}

func setRoute() {
	proxy := handler.NewGateway()
	recoveryMiddleware := middleware.NewRecover()
	upstreamMiddleware := middleware.NewUpstream()
	corsMiddleware := middleware.NewCors()
	rateLimitMiddleware := middleware.NewRateLimit(helper.NewRateLimiters())
	circuitBreakerMiddleware := middleware.NewCircuitBreaker(helper.NewNewCircuitBreakers())
	logMiddleware := middleware.NewLog()
	authenticateMiddleware := middleware.NewAuthenticate()
	http.Handle("/", getHandler(
		&proxy,
		&recoveryMiddleware,
		&upstreamMiddleware,
		&corsMiddleware,
		&rateLimitMiddleware,
		&circuitBreakerMiddleware,
		&logMiddleware,
		&authenticateMiddleware,
	))
	http.Handle("/metrics", promhttp.HandlerFor(helper.NewRegistry(), promhttp.HandlerOpts{}))
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if isShuttinDown.Load() {
			helper.JSONResponse(w, http.StatusServiceUnavailable, map[string]string{}, map[string]string{
				"msg": "Shutting down",
			})
			return
		}
		helper.JSONResponse(w, http.StatusOK, map[string]string{}, map[string]string{
			"msg": "OK",
		})
	})
}
