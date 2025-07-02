package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
)

func Test_HttpServerValidate(t *testing.T) {
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
	time.Sleep(1 * time.Second) // Даем серверу стартануть
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
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error reading response body: %s", err.Error())
		return
	}

	if len(body) == 0 {
		t.Errorf("response body is empty")
		return
	}

	fmt.Println("get body:")
	fmt.Println(string(body))

	var token models.SessionCreateResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		t.Errorf("error unmarshalling response: %s", err.Error())
		return
	}
	SessionServiceUrl = fmt.Sprintf("http://%s:%s/join", host, port)
	req2 := models.ValidateCodeReq{
		UserId: "1",
		Code:   token.SessionId,
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
	if user.SessionId != token.SessionId {
		t.Fatal("response body does not match session id")
	}

}
