package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"xxx/shared"
)

type MessageType string

const (
	MessageTypeAnswer      = MessageType("result")
	MessageTypeQuestion    = MessageType("question")
	MessageTypeLeaderboard = MessageType("leaderboard")
	MessageTypeAck         = MessageType("game_start")
)

// ClientMessage describes what we get from the user
type ClientMessage struct {
	Option int `json:"option"` // chosen answer index
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

	// ------ if Type is MessageTypeAck ------
	IsGameStarted bool `json:"isGameStarted,omitempty"` // broadcasting ack when first question triggered

	// ------ 'question payload' response to admin (triggered on next_question event); if Type is MessageTypeQuestion ------
	QuestionIdx     int             `json:"questionId,omitempty"`      // if Type is MessageTypeAnswer or MessageTypeQuestion
	QuestionsAmount int             `json:"questionsAmount,omitempty"` //
	Text            string          `json:"text,omitempty"`            // question text or feedback
	Options         []shared.Option `json:"options,omitempty"`         // for question

	// ------ if Type is MessageTypeAnswer ------
	Correct bool `json:"correct,omitempty"` // for answerResult

	//Payload     interface{}     `json:"payload,omitempty"`    // extra data (e.g. leaderboard)
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
func handleRead(ctx ConnectionContext, deps HandlerDeps) {
	defer func() {
		// On exit, clean up
		deps.Registry.UnregisterConnection(ctx.SessionId, ctx.UserId)
		ctx.Conn.Close()
	}()

	fmt.Println("Reach handleRead function for user", ctx.UserId)

	for {
		_, raw, err := ctx.Conn.ReadMessage()
		fmt.Printf("ws connection received message: '%s'\n", string(raw))
		if err != nil {
			fmt.Println(fmt.Errorf("ws error reading message: %w", err).Error())
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

		}
	}
}

// processAnswer processes an incoming UserMessage from a WebSocket client, then (optionally) sends immediate answer
func processAnswer(ctx ConnectionContext, deps HandlerDeps, msg *ClientMessage) {
	session := ctx.SessionId
	qid, _ := deps.Tracker.GetCurrentQuestion(ctx.SessionId)
	chosen := msg.Option

	// Look up the correct option from the QuizTracker
	correctIdx, correctOpt := deps.Tracker.GetCorrectOption(session, qid)
	if correctOpt == nil {
		log.Printf("no correct option found for session %s question %d", session, qid)
		return
	}

	// Record the answer
	deps.Tracker.RecordAnswer(session, ctx.UserId, chosen == correctIdx)

	// Send immediate feedback
	isCorrect := chosen == correctIdx
	resp := ServerMessage{
		Type:        MessageTypeAnswer,
		QuestionIdx: qid + 1,
		Correct:     isCorrect,
	}
	deps.Registry.sendMessage(resp.Bytes(), ctx)
}
