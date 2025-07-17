package Handlers

import (
	"log/slog"
	"xxx/LeaderBoardService/LeaderBoard"
)

type HandlerManager struct {
	log     *slog.Logger
	Service LeaderBoard.Service
}

func NewHandlerManager(log *slog.Logger, redisCon string) (*HandlerManager, error) {
	service, err := LeaderBoard.NewLeaderBoard(log, redisCon)
	if err != nil {
		return nil, err
	}
	return &HandlerManager{log: log, Service: service}, nil
}
