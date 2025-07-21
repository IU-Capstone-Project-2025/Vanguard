package HttpServer

import (
	"context"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xxx/LeaderBoardService/Handlers"
)

type HttpServer struct {
	*Handlers.HandlerManager
	Host     string
	Port     string
	logger   *slog.Logger
	server   *http.Server
	stopChan chan os.Signal
}

func InitHttpServer(logger *slog.Logger, Host string, Port string, RedisConn string) (*HttpServer, error) {
	logger.Info("InitHttpServer")
	managerHandler, err := Handlers.NewHandlerManager(logger, RedisConn)
	if err != nil {
		logger.Error("InitHttpServer", "NewSessionManagerHandler", err)
		return nil, err
	}
	return &HttpServer{
		HandlerManager: managerHandler,
		Host:           Host,
		Port:           Port,
		logger:         logger,
		stopChan:       make(chan os.Signal, 1),
	}, nil
}

func (hs *HttpServer) Start() {
	router := hs.registerHandlers()

	hs.server = &http.Server{
		Addr:    hs.Host + ":" + hs.Port,
		Handler: router,
	}

	// Захват SIGINT / SIGTERM
	signal.Notify(hs.stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		hs.logger.Info("HTTP server is starting", "addr", hs.server.Addr)
		if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hs.logger.Error("ListenAndServe error", "err", err)
		}
	}()

	<-hs.stopChan
	hs.logger.Info("Shutdown signal received")
	hs.Stop()
}

func (hs *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hs.logger.Info("Shutting down HTTP server...")
	if err := hs.server.Shutdown(ctx); err != nil {
		hs.logger.Error("HTTP server Shutdown", "err", err)
	} else {
		hs.logger.Info("HTTP server exited properly")
	}
}

// corsMiddleware is a middleware function that sets appropriate headers to http.ResponseWriter object
// to allow origins for CORS policy
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (hs *HttpServer) registerHandlers() *mux.Router {
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.HandleFunc("/get-results", hs.ComputeBoardHandler).Methods("POST", "OPTIONS")
	hs.logger.Info("Routes registered", "host", hs.Host, "port", hs.Port)
	return router
}
