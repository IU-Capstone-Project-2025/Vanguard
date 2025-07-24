package integration_tests

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/websocket"
//	"github.com/joho/godotenv"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"os"
//	"testing"
//	"time"
//	"xxx/SessionService/httpServer"
//	"xxx/SessionService/models"
//)
//
//func Test_HttpCloseWs(t *testing.T) {
//	cwd, _ := os.Getwd()
//	fmt.Println("Working dir:", cwd)
//
//	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
//		if err := godotenv.Load(getEnvFilePath()); err != nil {
//			t.Fatalf("could not load .env file: %v", err)
//		}
//	}
//	host := os.Getenv("SESSION_SERVICE_HOST")
//	port := os.Getenv("SESSION_SERVICE_PORT")
//
//	rabbitC, rabbitURL := startRabbit(context.Background(), t)
//	redisC, redisURL := startRedis(context.Background(), t)
//	defer redisC.Terminate(context.Background())
//	defer rabbitC.Terminate(context.Background())
//	log := setupLogger(envLocal)
//	server, err := httpServer.InitHttpServer(log, host, port, rabbitURL, redisURL)
//	if err != nil {
//		t.Fatalf("error creating http server: %v", err)
//	}
//	go server.Start()
//	time.Sleep(2 * time.Second)
//	defer server.Stop()
//
//	SessionServiceUrl := fmt.Sprintf("http://%s:%s/sessionsMock", host, port)
//	req := models.CreateSessionReq{
//		UserName: "admin",
//		QuizId:   "d2372184-dedf-42db-bcbd-d6bb15b0712b",
//	}
//	jsonBytes, err := json.Marshal(req)
//	if err != nil {
//		t.Fatal("error marshaling json:", err)
//	}
//
//	resp, err := http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
//	if err != nil {
//		t.Fatal("error making request:", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
//	}
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatalf("error reading response body: %v", err)
//	}
//
//	var token models.SessionCreateResponse
//	err = json.Unmarshal(body, &token)
//	if err != nil {
//		t.Fatalf("error unmarshalling response: %v", err)
//	}
//	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
//	req2 := models.ValidateCodeReq{
//		UserName: "user1",
//		Code:     token.SessionId,
//	}
//	jsonBytes, err = json.Marshal(req2)
//	if err != nil {
//		t.Fatal("error marshaling json:", err)
//	}
//
//	resp, err = http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
//	if err != nil {
//		t.Fatal("error making request:", err)
//	}
//	defer resp.Body.Close()
//	body2, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal("error reading response body:", err)
//	}
//	var user models.SessionCreateResponse
//	err = json.Unmarshal(body2, &user)
//	if err != nil {
//		t.Fatal("error unmarshalling response body:", err)
//	}
//
//	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
//	req3 := models.ValidateCodeReq{
//		UserName: "user2",
//		Code:     token.SessionId,
//	}
//	jsonBytes, err = json.Marshal(req3)
//	if err != nil {
//		t.Fatal("error marshaling json:", err)
//	}
//
//	resp, err = http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
//	if err != nil {
//		t.Fatal("error making request:", err)
//	}
//	defer resp.Body.Close()
//	body3, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal("error reading response body:", err)
//	}
//	var user2 models.SessionCreateResponse
//	err = json.Unmarshal(body3, &user2)
//	if err != nil {
//		t.Fatal("error unmarshalling response body:", err)
//	}
//
//	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
//	req4 := models.ValidateCodeReq{
//		UserName: "user3",
//		Code:     token.SessionId,
//	}
//	jsonBytes, err = json.Marshal(req4)
//	if err != nil {
//		t.Fatal("error marshaling json:", err)
//	}
//
//	resp, err = http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
//	if err != nil {
//		t.Fatal("error making request:", err)
//	}
//	defer resp.Body.Close()
//	body4, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal("error reading response body:", err)
//	}
//	var user3 models.SessionCreateResponse
//	err = json.Unmarshal(body4, &user3)
//	if err != nil {
//		t.Fatal("error unmarshalling response body:", err)
//	}
//
//	// Goroutine для user1
//	go func() {
//		u := url.URL{
//			Scheme:   "ws",
//			Host:     "localhost:8081",
//			Path:     "/ws",
//			RawQuery: "token=" + user.Jwt,
//		}
//		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//		if err != nil {
//			t.Error("user1 dial error:", err)
//			return
//		}
//		for i := 0; i < 3; i++ {
//			_, msg, err := conn.ReadMessage()
//			if err != nil {
//				t.Errorf("user1 read error: %v", err)
//				return
//			}
//			t.Log(string(msg))
//		}
//	}()
//
//	// Goroutine для user2
//	go func() {
//		time.Sleep(1 * time.Second)
//		u := url.URL{
//			Scheme:   "ws",
//			Host:     "localhost:8081",
//			Path:     "/ws",
//			RawQuery: "token=" + user2.Jwt,
//		}
//		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//		defer conn.Close()
//		if err != nil {
//			t.Error("user2 dial error:", err)
//			return
//		}
//		for i := 0; i < 2; i++ {
//			_, msg, err := conn.ReadMessage()
//			if err != nil {
//				t.Errorf("user2 read error: %v", err)
//				return
//			}
//			t.Log(string(msg))
//		}
//	}()
//
//	go func() {
//		time.Sleep(10 * time.Second)
//		u := url.URL{
//			Scheme:   "ws",
//			Host:     "localhost:8081",
//			Path:     "/ws",
//			RawQuery: "token=" + user3.Jwt,
//		}
//		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//		if err != nil {
//			t.Fatal("user3 dial error:", err)
//			return
//		}
//		_, msg, err := conn.ReadMessage()
//		if err != nil {
//			t.Fatalf("user2 read error: %v", err)
//			return
//		}
//		t.Log(string(msg))
//	}()
//	select {}
//}
