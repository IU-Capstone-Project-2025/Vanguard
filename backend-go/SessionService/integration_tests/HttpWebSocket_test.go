package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
)

func Test_HttpWebSocket(t *testing.T) {
	cwd, _ := os.Getwd()
	fmt.Println("Working dir:", cwd)

	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	host := os.Getenv("SESSION_SERVICE_HOST")
	port := os.Getenv("SESSION_SERVICE_PORT")

	rabbitC, rabbitURL := startRabbit(context.Background(), t)
	redisC, redisURL := startRedis(context.Background(), t)
	defer redisC.Terminate(context.Background())
	defer rabbitC.Terminate(context.Background())
	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, host, port, rabbitURL, redisURL)
	if err != nil {
		t.Fatalf("error creating http server: %v", err)
	}
	go server.Start()
	time.Sleep(2 * time.Second)
	defer server.Stop()

	SessionServiceUrl := fmt.Sprintf("http://%s:%s/sessionsMock", host, port)
	req := models.CreateSessionReq{
		UserId: "1",
		QuizId: "d2372184-dedf-42db-bcbd-d6bb15b0712b",
	}
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal("error marshaling json:", err)
	}

	resp, err := http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal("error making request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}

	var token models.SessionCreateResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
	req2 := models.ValidateCodeReq{
		UserId:   "test1",
		UserName: "user1",
		Code:     token.SessionId,
	}
	jsonBytes, err = json.Marshal(req2)
	if err != nil {
		t.Fatal("error marshaling json:", err)
	}

	resp, err = http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal("error making request:", err)
	}
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading response body:", err)
	}
	var user models.SessionCreateResponse
	err = json.Unmarshal(body2, &user)
	if err != nil {
		t.Fatal("error unmarshalling response body:", err)
	}

	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
	req3 := models.ValidateCodeReq{
		UserId:   "test2",
		UserName: "user2",
		Code:     token.SessionId,
	}
	jsonBytes, err = json.Marshal(req3)
	if err != nil {
		t.Fatal("error marshaling json:", err)
	}

	resp, err = http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal("error making request:", err)
	}
	defer resp.Body.Close()
	body3, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading response body:", err)
	}
	var user2 models.SessionCreateResponse
	err = json.Unmarshal(body3, &user2)
	if err != nil {
		t.Fatal("error unmarshalling response body:", err)
	}
	user1Chan := make(chan string, 2)
	user2Chan := make(chan string, 1)

	// Goroutine для user1
	go func() {
		u := url.URL{
			Scheme:   "ws",
			Host:     "localhost:8081",
			Path:     "/ws",
			RawQuery: "token=" + user.Jwt,
		}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			t.Fatal("user1 dial error:", err)
			return
		}
		defer conn.Close()

		for i := 0; i < 2; i++ {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				t.Errorf("user1 read error: %v", err)
				return
			}
			user1Chan <- string(msg)
		}
	}()

	// Goroutine для user2
	go func() {
		time.Sleep(1 * time.Second)
		u := url.URL{
			Scheme:   "ws",
			Host:     "localhost:8081",
			Path:     "/ws",
			RawQuery: "token=" + user2.Jwt,
		}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			t.Fatal("user2 dial error:", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("user2 read error: %v", err)
			return
		}
		user2Chan <- string(msg)
	}()

	// --- Сбор всех сообщений
	var (
		msg1a, msg1b, msg2 string
		received           int
		timeout            = time.After(5 * time.Second)
	)

	for received < 3 {
		select {
		case m := <-user1Chan:
			if msg1a == "" {
				msg1a = m
			} else {
				msg1b = m
			}
			received++
		case m := <-user2Chan:
			msg2 = m
			received++
		case <-timeout:
			t.Fatal("Timeout waiting for WebSocket messages")
		}
	}

	// --- Проверка содержимого
	expected1a := `["user1"]`
	expected1b := `["user2"]`
	expected2 := `["user2","user1"]`
	expected3 := `["user1","user2"]`

	if msg1a != expected1a {
		t.Fatalf("user1 first message mismatch. Got: %s, Want: %s", msg1a, expected1a)
	}
	if msg1b != expected1b {
		t.Fatalf("user1 second message mismatch. Got: %s, Want: %s", msg1b, expected1b)
	}
	if msg2 != expected2 && msg2 != expected3 {
		t.Fatalf("user2 message mismatch. Got: %s, Want: %s", msg2, expected2)
	}

}
