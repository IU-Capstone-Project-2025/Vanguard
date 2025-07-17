package ws

import (
	"fmt"
	"xxx/real_time/models"
	"xxx/shared"
)

type Responder struct {
	registry  *ConnectionRegistry
	sessionId string
}

func NewResponder(reg *ConnectionRegistry, sid string) Responder {
	return Responder{registry: reg, sessionId: sid}
}

func (r Responder) SendGameEnd() {
	gameEndAck := ServerMessage{
		Type: MessageTypeEnd,
	}
	r.registry.BroadcastToSession(r.sessionId, gameEndAck.Bytes(), false)
}

func (r Responder) SendNextQuestionAck() {
	nextQuestionAck := ServerMessage{
		Type: MessageTypeNextQuestion,
	}
	r.registry.BroadcastToSession(r.sessionId, nextQuestionAck.Bytes(), false)

	fmt.Printf("Send next question ack for %sid: %v\n", r.sessionId, nextQuestionAck)
}

func (r Responder) SendError() {
	gameEndAck := ServerMessage{
		Type: MessageTypeError,
	}
	r.registry.BroadcastToSession(r.sessionId, gameEndAck.Bytes(), false)
}

func (r Responder) SendLeaderboard(lb shared.ScoreTable) {
	leaderBoard := ServerMessage{
		Type:    MessageTypeLeaderboard,
		Payload: lb,
	}

	r.registry.SendToAdmin(r.sessionId, leaderBoard.Bytes())
}

func (r Responder) SendQuestionStat(questionStat shared.PopularAns, questionAnswers map[string]models.UserAnswer) {
	for _, connectionCtx := range r.registry.GetConnections(r.sessionId) { // iterate through all connections to retrieve userIds
		if connectionCtx.Role == shared.RoleAdmin { // skip admin, since we do not send stat to him
			continue
		}
		user := connectionCtx.UserId

		stat := ServerMessage{
			Type:    MessageTypeStat,
			Correct: questionAnswers[user].Correct,
			Payload: questionStat,
		}

		r.registry.SendMessage(stat.Bytes(), connectionCtx)
	}
}

func (r Responder) SendQuestionPayload(qid, questionsAmount int, question shared.Question) {
	questionPayloadMsg := ServerMessage{
		Type:            MessageTypeQuestion,
		QuestionIdx:     qid,
		QuestionsAmount: questionsAmount,
		Text:            question.Text,
		Options:         question.Options,
		Payload:         question.ImageUrl,
	}

	r.registry.SendToAdmin(r.sessionId, questionPayloadMsg.Bytes())
	fmt.Printf("Send question payload for %sid: %v\n", r.sessionId, questionPayloadMsg)
}
