package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"xxx/real_time/ws"
)

func ReadWs(t *testing.T, conn *websocket.Conn) ws.ServerMessage {
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var serverMsg ws.ServerMessage
	err = json.Unmarshal(msg, &serverMsg)

	t.Logf("Received: %s", msg)
	return serverMsg
}

func ConnectWs(t *testing.T, token string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}

	return conn
}

func CloseWs(conn *websocket.Conn) {
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
	conn.Close()
}
