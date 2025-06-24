package shared

const (
	SessionExchange = "session.events" // the name of the only exchange

	SessionStartRoutingKey = "session.start" // routing key for "session_start" event
	SessionEndRoutingKey   = "session.end"   // routing key for "session_end" event
	// QuestionStartRoutingKey is a routing key for "question_start" event.
	QuestionStartRoutingKey = "question.*.start" // * stands for session code
)
