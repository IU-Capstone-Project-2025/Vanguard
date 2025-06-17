package Handlers

import "xxx/SessionService/Game"

type SessionManagerHandler struct {
	Manager Game.Manager
}

func NewSessionManagerHandler(manager Game.Manager) *SessionManagerHandler {
	return &SessionManagerHandler{Manager: manager}
}
