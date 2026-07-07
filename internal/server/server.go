package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/evgen6501-star/golang-subscriptions/internal/config"
	"github.com/evgen6501-star/golang-subscriptions/internal/logger"
	"github.com/evgen6501-star/golang-subscriptions/internal/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     config.ServerConfig
	log        *logger.Logger
	middleware []middleware.Middleware
}

func NewHttpServer(config config.ServerConfig, log *logger.Logger, middle ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middle,
	}
}
func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}
func (h *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(h.mux, h.middleware...)
	server := &http.Server{
		Addr:    h.config.Host + ":" + h.config.Port,
		Handler: mux,
	}
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		h.log.Warn("start HTTP server", zap.String("addr", h.config.Host))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()
	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listend and sever HTTP :%w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shotdown HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5,
		)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shotdown HTTP server :%w", err)
		}
		h.log.Warn("HTTP server stopped")
	}
	return nil
}
