package httpServer

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"xxx/SessionService/Handlers"
	_ "xxx/SessionService/docs"
)

type HttpServer struct {
	*Handlers.SessionManagerHandler
	Host   string
	Port   string
	logger *slog.Logger
}

func InitHttpServer(logger *slog.Logger, Host string, Port string, rmqConn string, RedisConn string) (*HttpServer, error) {
	logger.Info("InitHttpServer")
	managerHandler, err := Handlers.NewSessionManagerHandler(rmqConn, RedisConn, logger)
	if err != nil {
		logger.Error("InitHttpServer", "NewSessionManagerHandler", err)
		return nil, err
	}
	return &HttpServer{managerHandler, Host, Port, logger}, nil
}

func (HttpServer *HttpServer) Start() {
	go func() {
		err := HttpServer.registerHandlers()
		if err != nil {
			HttpServer.logger.Error("StartHttpServer", "RegisterHandlers", err)
			return
		}
		HttpServer.logger.Info("StartHttpServer", "StartHttpServer ok")
	}()
	select {}
}

func (HttpServer *HttpServer) Stop() {
	//TODO implement function
}

func (HttpServer *HttpServer) registerHandlers() error {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/create", HttpServer.CreateSessionHandler).Methods("GET")
	router.HandleFunc("/validate", HttpServer.ValidateCodeHandler).Methods("GET")
	router.HandleFunc("/next", HttpServer.NextQuestionHandler).Methods("GET")
	router.HandleFunc("/start", HttpServer.StartSessionHandler).Methods("GET")
	HttpServer.logger.Info("registerHandlers", "msg", "Listening on "+HttpServer.Host+":"+HttpServer.Port)
	err := http.ListenAndServe(HttpServer.Host+":"+HttpServer.Port, router)
	if err != nil {
		HttpServer.logger.Error("registerHandlers", "ListenAndServe", err)
		return err
	}
	return nil
}
