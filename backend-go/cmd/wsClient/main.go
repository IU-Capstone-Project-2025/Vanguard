package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <jwt-token>")
		return
	}

	token := os.Args[1]
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid JWT format")
		return
	}

	payloadPart := parts[1]

	// Добавим padding если нужно
	switch len(payloadPart) % 4 {
	case 2:
		payloadPart += "=="
	case 3:
		payloadPart += "="
	case 1:
		payloadPart += "==="
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(payloadPart)
	if err != nil {
		fmt.Println("Error decoding payload:", err)
		return
	}

	var payload map[string]interface{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		fmt.Println("Error parsing JSON payload:", err)
		return
	}

	// Выводим расшифрованный payload
	fmt.Println("Decoded JWT payload:")
	for k, v := range payload {
		fmt.Printf("%s: %v\n", k, v)
	}
	scheme := "ws"           // или "wss" если HTTPS
	host := "localhost:8081" // твой сервер и порт
	path := "/ws"
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	fmt.Printf("Connecting to %s\n", u.String())

	// Подключаемся
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	// Перехватываем Ctrl+C, чтобы корректно закрыть соединение
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Горутина для чтения сообщений с сервера
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	// В этом простом примере мы просто ждём прерывания Ctrl+C
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("interrupt received, closing connection")
			// Отправляем close сообщение серверу
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
				return
			}
			return
		}
	}
}
