package shared

// Session represents a message pushed to RabbitMQ as "session_"-type event
type Session struct {
	ID               string `redis:"id"`
	Code             string `redis:"code"`
	State            string `redis:"state"`
	ServerWsEndpoint string `redis:"serverWsEndpoint"`
}
