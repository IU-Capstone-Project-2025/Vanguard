package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"xxx/real_time/models"
)

// handleRead continuously reads messages from the WebSocket connection.
// It gets incoming messages and delegates processing to HandleUserMessage.
// If an error occurs (e.g., due to a disconnect), it ensures the connection is closed gracefully.
func handleRead(ctx *ConnectionContext, reg *ConnectionRegistry) {
	defer func() {
		// On exit, clean up
		reg.UnregisterConnection(ctx.SessionId, ctx.UserId)
		ctx.Conn.Close()
	}()

	err := reg.RegisterConnection(ctx.SessionId, ctx.UserId, ctx.Conn)
	if err != nil {
		fmt.Println("ws error in token: %w", err)
		return
	}

	for {
		_, msg, err := ctx.Conn.ReadMessage()
		if err != nil {
			fmt.Println("ws error reading message: %w", err)
			return
		}

		handleUserMessage(msg, ctx)
	}
}

// handleUserMessage processes an incoming UserMessage from a WebSocket client.
// TODO: implement real functionality
func handleUserMessage(msg []byte, ctx *ConnectionContext) {
	reply := ""
	if ctx.Role == models.RoleAdmin {
		reply = "Hi, room's Host. Echo:\n" + string(msg)
	} else {
		reply = "Hi, user. Echo:\n" + string(msg)
	}

	err := ctx.Conn.WriteMessage(websocket.TextMessage, []byte(reply))
	if err != nil {
		fmt.Println("ws error writing message:", err)
	}
}
