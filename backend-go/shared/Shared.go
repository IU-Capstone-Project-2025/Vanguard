package shared

const WsEndpoint string = "/ws/server"
const CodeLength int = 6
const (
	SessionExchange         = "session.events"
	SessionStartRoutingKey  = "session.start"
	SessionEndRoutingKey    = "session.end"
	QuestionStartRoutingKey = "question.*.start"
)
