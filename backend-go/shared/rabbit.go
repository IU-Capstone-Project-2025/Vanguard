package shared

const (
	SessionExchange = "session.events"

	SessionStartRoutingKey  = "session.start"
	SessionEndRoutingKey    = "session.end"
	QuestionStartRoutingKey = "question.*.start"
)

type RabbitSessionMsg struct {
	SessionId string
}
