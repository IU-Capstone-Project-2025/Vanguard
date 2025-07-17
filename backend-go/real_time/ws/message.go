package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"xxx/real_time/models"
	"xxx/shared"
)

type MessageType string

const (
	MessageTypeAnswer   = MessageType("answer")
	MessageTypeQuestion = MessageType("question")

	MessageTypeLeaderboard = MessageType("leaderboard")
	MessageTypeStat        = MessageType("question_stat")

	MessageTypeEnd = MessageType("end") // sent to admin when game ends

	MessageTypeNextQuestion = MessageType("next_question") // sent to admin when next question is triggered
	MessageTypeUserAnswered = MessageType("user_answered") // sent to admin when participant submitted his answer

	MessageTypeError = MessageType("error")
)

// ClientMessage describes what we get from the user
type ClientMessage struct {
	Type MessageType `json:"type"` // MessageTypeAnswer or MessageTypeNextQuestion

	// ------ if Type is MessageTypeAnswer ------
	Option    int       `json:"option,omitempty"`    // chosen answer index
	Timestamp time.Time `json:"timestamp,omitempty"` // time user have answered
}

func (m *ClientMessage) Bytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("failed to marshal ClientMessage: %v", err)
	}
	return b
}

// ServerMessage describes what we send back (to quiz host or participant).
type ServerMessage struct {
	// ------ response message type; on each response ------
	Type MessageType `json:"type"` // "question", "result", "leaderboard"

	// ------ 'question payload' response to admin (triggered on next_question event); if Type is MessageTypeQuestion ------
	QuestionIdx     int             `json:"questionId,omitempty"`      // if Type is MessageTypeAnswer or MessageTypeQuestion
	QuestionsAmount int             `json:"questionsAmount,omitempty"` //
	Text            string          `json:"text,omitempty"`            // question text or feedback
	Options         []shared.Option `json:"options,omitempty"`         // for question

	// ------ if Type is MessageTypeAnswer or MessageTypeStat ------
	Correct bool `json:"correct,omitempty"` // for answerResult

	// ------ if Type is MessageTypeLeaderboard or MessageTypeStat ------
	Payload interface{} `json:"payload,omitempty"` // extra data (e.g. leaderboard)
}

func (m *ServerMessage) Bytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("failed to marshal ServerMessage: %v", err)
	}
	return b
}

// handleRead continuously reads messages from the WebSocket connection.
// It gets incoming messages and delegates processing to HandleUserMessage.
// If an error occurs (e.g., due to a disconnect), it ensures the connection is closed gracefully.
func handleRead(ctx *ConnectionContext, deps HandlerDeps) {
	defer func() {
		// On exit, clean up
		deps.Registry.UnregisterConnection(ctx.SessionId, ctx.UserId)
		fmt.Println("CLOSING GA")
		ctx.Conn.Close()
	}()

	fmt.Println("Reach handleRead function for user", ctx.UserId)

	for {
		_, raw, err := ctx.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println(fmt.Errorf("ws error reading message: %w", err).Error())
			} else {
				log.Printf("websocket closed: %v\n", err)
			}
			return
		}

		var msg ClientMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			fmt.Printf("invalid ws message: %v", err)
			continue
		}

		switch ctx.Role {
		case shared.RoleParticipant:
			go processAnswer(ctx, deps, &msg)
		case shared.RoleAdmin:
			go func() {
				responder := NewResponder(deps.Registry, ctx.SessionId)
				responder.SendNextQuestionAck()
			}()
		}
	}
}

// processAnswer processes an incoming UserMessage from a WebSocket client, then (optionally) sends immediate answer
func processAnswer(ctx *ConnectionContext, deps HandlerDeps, msg *ClientMessage) {
	sessionId := ctx.SessionId
	qid, _ := deps.Tracker.GetCurrentQuestion(ctx.SessionId)

	// Look up the correct option from the QuizTracker
	correctIdx, correctOpt := deps.Tracker.GetCorrectOption(sessionId, qid)
	if correctOpt == nil {
		log.Printf("no correct option found for sessionId %s question %d", sessionId, qid)
		return
	}

	isCorrect := msg.Option == correctIdx
	userAnswer := models.UserAnswer{
		Option:    msg.Option,
		Answered:  true,
		Correct:   isCorrect,
		Timestamp: msg.Timestamp,
	}

	// Record the answer
	deps.Tracker.RecordAnswer(sessionId, ctx.UserId, userAnswer)
	fmt.Println("recorded answer ", userAnswer, "from ", ctx.UserId)

	// notify admin about new answered user
	resp := ServerMessage{
		Type: MessageTypeUserAnswered,
		Payload: map[string]string{
			"userId": ctx.UserId,
		},
	}
	deps.Registry.SendToAdmin(sessionId, resp.Bytes())
}
