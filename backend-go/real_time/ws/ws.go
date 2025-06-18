package ws

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"net/http"
	"xxx/real_time/models"
	"xxx/shared"
)

var Connections = make(map[string]map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  5 * 1024,
	WriteBufferSize: 5 * 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("tokenString")
	if tokenString == "" {

	}

	token, err := extractTokenData(tokenString)
	if err != nil {
		fmt.Println(err)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	ctx := &models.ConnectionContext{
		Conn:      conn,
		UserId:    token.UserId,
		SessionId: token.SessionId,
		Role:      token.Role,
	}

	go handleRead(conn, ctx)
}

func extractTokenData(tokenString string) (*shared.UserTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shared.UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(shared.JwtKey), nil
	}, nil)

	switch {
	case token.Valid:
		fmt.Println("OK token")
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("Malformed token")
	default:
		fmt.Println("Couldn't handle this token:", err)
	}

	claims, ok := token.Claims.(*shared.UserTokenClaims)
	if !ok {
		return nil, fmt.Errorf("error decoding jwt")
	}

	return claims, nil
}

func handleRead(conn *websocket.Conn, ctx *models.ConnectionContext) {
	defer conn.Close()

	registerConnection(ctx.SessionId, ctx.UserId, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Errorf("ws error reading message: %w", err)
			return
		}

		handleUserMessage(msg)
	}
}

func handleUserMessage(msg []byte) {

}

func registerConnection(sessionId string, userId string, conn *websocket.Conn) {
	Connections[sessionId][userId] = conn
}

func unregisterConnection(sessionId string, userId string) {
	delete(Connections[sessionId], userId)
}

func RegisterSession(msg shared.RabbitSessionMsg) {
	Connections[msg.SessionId] = make(map[string]*websocket.Conn)
}
