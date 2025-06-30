package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"xxx/shared"
)

type MessageType string

const (
	MessageTypeAnswer      = MessageType("result")
	MessageTypeQuestion    = MessageType("question")
	MessageTypeLeaderboard = MessageType("leaderboard")
)

// ClientMessage describes what we get from the user
type ClientMessage struct {
	Option int `json:"option"` // chosen answer index
}

// ServerMessage describes what we send back (to quiz host or participant).
type ServerMessage struct {
	Type       string          `json:"type"`                 // "question", "result", "leaderboard"
	QuestionID string          `json:"questionId,omitempty"` // for question & answerResult
	Text       string          `json:"text,omitempty"`       // question text or feedback
	Options    []shared.Option `json:"options,omitempty"`    // for question
	Correct    *bool           `json:"correct,omitempty"`    // for answerResult
	Payload    interface{}     `json:"payload,omitempty"`    // extra data (e.g. leaderboard)
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
		ctx.Conn.Close()
	}()

	err := deps.Registry.RegisterConnection(ctx.SessionId, ctx.UserId, ctx.Conn)
	if err != nil {
		fmt.Printf("failed to register joined user: %v", err)
		return
	}

	for {
		_, raw, err := ctx.Conn.ReadMessage()
		fmt.Printf("ws connection received message: '%s'\n", string(raw))
		if err != nil {
			fmt.Println(fmt.Errorf("ws error reading message: %w", err).Error())
			return
		}

		var msg ClientMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("invalid ws message: %v", err)
			continue
		}

		switch ctx.Role {
		case shared.RoleParticipant:
			go processAnswer(ctx, deps, &msg)
		case shared.RoleAdmin:
			log.Printf("ws message from the quizz host is ignored")
		}
	}
}

// processAnswer processes an incoming UserMessage from a WebSocket client, then (optionally) sends immediate answer
func processAnswer(ctx *ConnectionContext, deps HandlerDeps, msg *ClientMessage) {
	//session := ctx.SessionId
	//qid, _ := deps.Tracker.GetCurrentQuestionIdx(ctx.SessionId)
	//chosen := msg.Option
	//
	//// Look up the correct option from the QuizTracker
	//correctOpt, ok := deps.Tracker.GetCurrentQuestionIdx(session)
	//if !ok {
	//	log.Printf("no correct option found for session %s question %d", session, qid)
	//	return
	//}

	// Record the answer
	//deps.Tracker.RecordAnswer(session, ctx.UserId, qid, chosen, chosen == correctOpt)
	//
	//// Send immediate feedback
	//isCorrect := (chosen == correctOpt)
	//resp := ServerMessage{
	//	Type:       "answerResult",
	//	QuestionID: qid,
	//	Correct:    &isCorrect,
	//}
	//deps.Registry.SendToConn(ctx.Conn, resp.Bytes())
}

// sendMessage
func sendMessage(payload []byte, receivers ...*websocket.Conn) {
	for _, r := range receivers {
		err := r.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Printf("Failed to send message to %v: %v", r, err)
		}
		log.Println("message sent successfully")
	}
}
