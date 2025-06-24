package Handlers

import (
	"log/slog"
	"xxx/SessionService/Game"
	"xxx/shared"
)

type SessionManagerHandler struct {
	Manager Game.Manager
	logger  *slog.Logger
}

func NewSessionManagerHandler(rmqConn string, RedisConn string, log *slog.Logger) (*SessionManagerHandler, error) {
	codelenhtg := shared.CodeLength
	manager, err := Game.CreateSessionManager(codelenhtg, rmqConn, RedisConn)
	if err != nil {
		log.Error("error on NewSessionManagerHandler", err)
		return nil, err
	}
	return &SessionManagerHandler{Manager: manager, logger: log}, nil
}
