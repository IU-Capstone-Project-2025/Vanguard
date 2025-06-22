package models

// Role represent a user type
type Role string

const (
	RoleAdmin       = "admin"       // the host of the session (quiz)
	RoleParticipant = "participant" // the participant of the session (quiz)
)
