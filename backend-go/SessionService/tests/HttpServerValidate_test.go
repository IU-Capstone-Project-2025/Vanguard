package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
)

func Test_HttpServerValidate(t *testing.T) {
	log := setupLogger(envLocal)

	server, err := httpServer.InitHttpServer(log, "localhost", "8000", "amqp://guest:guest@localhost:5672/", "localhost:6379")
	if err != nil {
		log.Error("error creating http server", "error", err)
		return
	}

	go server.Start()
	time.Sleep(1 * time.Second) // Даем серверу стартануть

	resp, err := http.Get("http://localhost:8000/create?userId=1")
	if err != nil {
		t.Error("error making request:", err)
		return
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

	var token models.UserToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		t.Errorf("error unmarshalling response: %s", err.Error())
		return
	}
	u, err := url.Parse("http://localhost:8000/validate")
	if err != nil {
		t.Error("error parsing url:", err)
	}
	params := url.Values{}
	params.Add("userId", "1")
	params.Add("code", token.CurrentQuiz)
	u.RawQuery = params.Encode()
	resp2, err := http.Get(u.String())
	if err != nil {
		t.Error("error making request:", err)
	}
	defer resp2.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Error("error reading response body:", err)
	}
	var user models.UserToken
	err = json.Unmarshal(body2, &user)
	if err != nil {
		t.Error("error unmarshalling response body:", err)
	}
	fmt.Println("get body:", user)

}
