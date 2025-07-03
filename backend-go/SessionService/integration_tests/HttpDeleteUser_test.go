package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
)

func Test_HttpDeleteUser(t *testing.T) {
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

	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
	req4 := models.ValidateCodeReq{
		UserId:   "test3",
		UserName: "user3",
		Code:     token.SessionId,
	}
	jsonBytes, err = json.Marshal(req4)
	if err != nil {
		t.Fatal("error marshaling json:", err)
	}

	resp, err = http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal("error making request:", err)
	}
	defer resp.Body.Close()
	body4, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading response body:", err)
	}
	var user3 models.SessionCreateResponse
	err = json.Unmarshal(body4, &user3)
	if err != nil {
		t.Fatal("error unmarshalling response body:", err)
	}
	user1Chan := make(chan string, 2)
	user2Chan := make(chan string, 1)
	user3Chan := make(chan string, 1)

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
	expected1a := `{"test1":"user1"}`
	expected1b := `{"test2":"user2"}`
	expected2 := `{"test1":"user1","test2":"user2"}`
	expected3 := `{"test2":"user2","test1":"user1"}`

	if msg1a != expected1a {
		t.Fatalf("user1 first message mismatch. Got: %s, Want: %s", msg1a, expected1a)
	}
	if msg1b != expected1b {
		t.Fatalf("user1 second message mismatch. Got: %s, Want: %s", msg1b, expected1b)
	}
	if msg2 != expected2 && msg2 != expected3 {
		t.Fatalf("user2 message mismatch. Got: %s, Want: %s", msg2, expected2)
	}
	uuuu := fmt.Sprintf("http://%s:%s/delete-user?code=%s&userId=%s", host, port, token.SessionId, "test1")
	resp, err = http.Get(uuuu)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Request failed: %v", resp.StatusCode)
	}

	go func() {
		time.Sleep(1 * time.Second)
		u := url.URL{
			Scheme:   "ws",
			Host:     "localhost:8081",
			Path:     "/ws",
			RawQuery: "token=" + user3.Jwt,
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
		fmt.Println(string(msg))
		user3Chan <- string(msg)

	}()
	timeout2 := time.After(5 * time.Second)
	var msg3 string
	select {
	case m := <-user3Chan:
		msg3 = m
	case <-timeout2:
		t.Fatal("Timeout waiting for WebSocket messages")
	}
	expected4 := `{"test2":"user2","test3":"user3"}`
	expected5 := `{"test3":"user3","test2":"user2"}`
	if msg3 != expected4 && msg3 != expected5 {
		t.Fatalf("user3 second message mismatch. Got: %s, Want: %s", msg3, expected5)
	}
}
