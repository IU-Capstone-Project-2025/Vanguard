package integration_tests

import (
	"bytes"
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

	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	host := os.Getenv("SESSION_SERVICE_HOST")
	port := os.Getenv("SESSION_SERVICE_PORT")

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))

	redisURL := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, host, port, rabbitURL, redisURL)
	if err != nil {
		t.Fatalf("error creating http server: %v", err)
	}
	go server.Start()
	time.Sleep(1 * time.Second)

	SessionServiceUrl := fmt.Sprintf("http://%s:%s/sessions", host, port)
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
		t.Error("error reading response body:", err)
	}
	var user models.SessionCreateResponse
	err = json.Unmarshal(body2, &user)
	if err != nil {
		t.Error("error unmarshalling response body:", err)
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
		t.Error("error reading response body:", err)
	}
	var user2 models.SessionCreateResponse
	err = json.Unmarshal(body3, &user2)
	if err != nil {
		t.Error("error unmarshalling response body:", err)
	}
	go func() {
		scheme := "ws"             // или "wss" если HTTPS
		wsHost := "localhost:8081" // твой сервер и порт
		path := "/ws"
		u := url.URL{Scheme: scheme, Host: wsHost, Path: path}
		q := u.Query()
		q.Set("token", user.Jwt)
		u.RawQuery = q.Encode()

		fmt.Printf("Connecting to %s\n", u.String())

		// Подключаемся
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("dial error:", err)
		}
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			var m interface{}
			err = json.Unmarshal(message, &m)
			if err != nil {
				fmt.Println("read error:", err)
			}
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			fmt.Printf("Received user1: %s\n", m)
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		scheme := "ws"             // или "wss" если HTTPS
		wsHost := "localhost:8081" // твой сервер и порт
		path := "/ws"
		u := url.URL{Scheme: scheme, Host: wsHost, Path: path}
		q := u.Query()
		q.Set("token", user2.Jwt)
		u.RawQuery = q.Encode()

		fmt.Printf("Connecting to %s\n", u.String())

		// Подключаемся
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("dial error:", err)
		}
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			var m interface{}
			err = json.Unmarshal(message, &m)
			if err != nil {
				fmt.Println("read error:", err)
			}
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			fmt.Printf("Received user2: %s\n", m)
		}
	}()
	select {}

}
