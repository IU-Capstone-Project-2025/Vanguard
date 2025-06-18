package Shared

const WsEndpoint string = "/ws/hui"
const CodeLength int = 6
const (
	SessionExchange         = "session.events"
	SessionStartRoutingKey  = "session.start"
	SessionEndRoutingKey    = "session.end"
	QuestionStartRoutingKey = "question.*.start"
)
