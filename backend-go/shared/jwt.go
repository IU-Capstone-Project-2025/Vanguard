package shared

import (
	"github.com/golang-jwt/jwt/v5"
)

// UserRole represent a user type
type UserRole string

const (
	RoleAdmin       = "admin"       // the host of the session (quiz)
	RoleParticipant = "participant" // the participant of the session (quiz)
)

// UserToken represents the structure of the user's ephemeral token
type UserToken struct {
	UserId    string   `json:"userId"`
	UserType  UserRole `json:"userType"`
	SessionId string   `json:"sessionId"`
	Exp       int64    `json:"exp"`
	jwt.RegisteredClaims
}
