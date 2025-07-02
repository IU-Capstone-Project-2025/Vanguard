package httpServer

import (
	"context"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xxx/SessionService/Handlers"
	_ "xxx/SessionService/docs"
)

type HttpServer struct {
	*Handlers.SessionManagerHandler
	Host     string
	Port     string
	logger   *slog.Logger
	server   *http.Server
	stopChan chan os.Signal
}

func InitHttpServer(logger *slog.Logger, Host string, Port string, rmqConn string, RedisConn string) (*HttpServer, error) {
	logger.Info("InitHttpServer")
	managerHandler, err := Handlers.NewSessionManagerHandler(rmqConn, RedisConn, logger)
	if err != nil {
		logger.Error("InitHttpServer", "NewSessionManagerHandler", err)
		return nil, err
	}
	return &HttpServer{
		SessionManagerHandler: managerHandler,
		Host:                  Host,
		Port:                  Port,
		logger:                logger,
		stopChan:              make(chan os.Signal, 1),
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

func (hs *HttpServer) registerHandlers() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/sessions", hs.CreateSessionHandler).Methods("POST")
	router.HandleFunc("/join", hs.ValidateCodeHandler).Methods("POST")
	router.HandleFunc("/session/{id}/nextQuestion", hs.NextQuestionHandler).Methods("POST")
	router.HandleFunc("/start", hs.StartSessionHandler).Methods("POST")
	router.HandleFunc("/validate", hs.ValidateSessionCodeHandler).Methods("POST")
	router.HandleFunc("/sessionsMock", hs.CreateSessionHandlerMock).Methods("POST")
	router.HandleFunc("/session/{id}/end", hs.SessionEndHandler).Methods("POST")
	router.HandleFunc("/healthz", hs.HealthHandler).Methods("POST")
	registry := Handlers.NewConnectionRegistry(hs.logger)
	router.Handle("/ws", Handlers.NewWebSocketHandler(registry))
	router.Handle("/delete-user", Handlers.DeleteUserHandler(registry))
	hs.logger.Info("Routes registered", "host", hs.Host, "port", hs.Port)
	return router
}
