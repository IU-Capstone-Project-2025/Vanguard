package models

import "github.com/gorilla/websocket"

type Role string

const (
	RoleAdmin       = "admin"
	RoleParticipant = "participant"
)

type ConnectionContext struct {
	Conn      *websocket.Conn
	UserId    string
	SessionId string
	Role      Role
}
