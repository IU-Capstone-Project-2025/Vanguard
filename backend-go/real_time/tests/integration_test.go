//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/real_time/app"
	"xxx/real_time/config"
	"xxx/real_time/rabbit"
	"xxx/real_time/ws"
	"xxx/shared"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startRabbit(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	defenitionsAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "rabbit", "definitions.json"))
	require.NoError(t, err)
	confAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "rabbit", "rabbitmq.conf"))
	require.NoError(t, err)

	// 1. Start RabbitMQ container
	rabbitReq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672:5672/tcp", "15672:15672/tcp"},
		Env: map[string]string{
			"RABBITMQ_LOAD_DEFINITIONS": "true",
			"RABBITMQ_DEFINITIONS_FILE": "/etc/rabbitmq/definitions.json",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      defenitionsAbs, // will be discarded internally
				ContainerFilePath: "/etc/rabbitmq/definitions.json",
				FileMode:          644,
			},

			{
				HostFilePath:      confAbs, // will be discarded internally
				ContainerFilePath: "/etc/rabbitmq/rabbitmq.conf",
				FileMode:          644,
			},
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}
	rabbitC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: rabbitReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start RabbitMQ container: %v", err)
	}

	rabbitHost, err := rabbitC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rabbitPort, err := rabbitC.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatal(err)
	}

	u := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.LoadConfig().MQ.User, config.LoadConfig().MQ.Password,
		rabbitHost, rabbitPort.Port())
	t.Logf("Rabbit running at %s", u)
	return rabbitC, u
}

func TestWithTestContainers(t *testing.T) {
	ctx := context.Background()

	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		fmt.Println("LOADING .ENV")
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}

	_, amqpURL := startRabbit(ctx, t)

	t.Log("rabbit running on ", amqpURL)

	// 2. Start Redis container similarly if needed

	// 3. Start RealTime service in a goroutine or exec.Command, configuring it to connect to amqpURL and Redis.
	//    For brevity, assume RealTime service can be started in-process or as a subprocess, reading env vars:
	// Start your RealTime main in a goroutine if possible, or exec binary.
	go startRealTimeServer(t, amqpURL)
	time.Sleep(2 * time.Second)

	adminId := "admin"
	users := []string{"navalniy"}
	sessionId := "В4ФЛ3Р"
	quiz := shared.Quiz{Questions: []shared.Question{
		{
			Type: "single_choice",
			Text: "What is the output of print(2 ** 3)?",
			Options: []shared.Option{
				{Text: "6", IsCorrect: false},
				{Text: "8", IsCorrect: true},
				{Text: "9", IsCorrect: false},
				{Text: "5", IsCorrect: false},
			},
		},
		{
			Type: "single_choice",
			Text: "Which keyword is used to create a function in Python?",
			Options: []shared.Option{
				{Text: "func", IsCorrect: false},
				{Text: "function", IsCorrect: false},
				{Text: "def", IsCorrect: true},
				{Text: "define", IsCorrect: false},
			},
		},
		{
			Type: "single_choice",
			Text: "What data type is the result of: 3 / 2 in Python 3?",
			Options: []shared.Option{
				{Text: "int", IsCorrect: false},
				{Text: "float", IsCorrect: true},
				{Text: "str", IsCorrect: false},
				{Text: "decimal", IsCorrect: false},
			},
		},
	}}

	usersAnswers := [][]int{
		{2, 2, 3},
	}

	adminToken := generateJWT(t, sessionId, adminId, shared.RoleAdmin)

	var adminConn *websocket.Conn
	var usersConn []*websocket.Conn

	// 5. Publish session.start
	publishSessionStart(t, amqpURL, sessionId, quiz)

	// 6. Admin WS connection
	adminConn = connectWs(t, adminToken)

	// 7. Read welcome
	readWs(t, adminConn)

	// join users
	for _, userId := range users {
		userToken := generateJWT(t, sessionId, userId, shared.RoleParticipant)
		conn := connectWs(t, userToken)
		usersConn = append(usersConn, conn)

		readWs(t, conn)
	}

	// 8. Start question flow
	for i, q := range quiz.Questions {
		t.Log("trigger question ", i, q)
		publishQuestionStart(t, amqpURL, sessionId)

		questionPayload := readWs(t, adminConn)
		t.Log("checking question payload:")
		require.Equal(t, q.Text, questionPayload.Text)
		require.Equal(t, ws.MessageTypeQuestion, questionPayload.Type)
		require.Equal(t, i+1, questionPayload.QuestionIdx)
		require.Equal(t, q.Options, questionPayload.Options)

		if i == 0 {
			for _, user := range usersConn {
				readWs(t, user)
			}
		}

		for j, user := range usersConn {
			option := usersAnswers[j][i]
			t.Logf("user %s send answer: %d", users[j], option)
			msg := ws.ClientMessage{Option: option}
			user.WriteMessage(websocket.TextMessage, msg.Bytes())

			resp := readWs(t, user)
			require.Equal(t, ws.MessageTypeAnswer, resp.Type)
			require.Equal(t, i+1, resp.QuestionIdx)
			require.Equal(t, q.Options[option].IsCorrect, resp.Correct)
		}
	}

	// trigger session end
	publishSessionEnd(t, amqpURL, sessionId)
	t.Log("---- Admin received leaderboard:")
	readWs(t, adminConn)

	t.Log("---- Users received leaderboards:")
	for i, user := range usersConn {
		lb := readWs(t, user)
		t.Log(lb.Payload)
		ans, ok := lb.Payload.(map[string]interface{})
		require.Equal(t, true, ok)

		userChosen, ok := ans[users[i]].([]interface{})
		require.Equal(t, true, ok)

		for j, isCorrectInter := range userChosen {
			chosenIdx := usersAnswers[i][j]

			isCorrect, ok := isCorrectInter.(bool)
			require.Equal(t, true, ok)

			require.Equal(t, quiz.Questions[j].Options[chosenIdx].IsCorrect, isCorrect)
		}
	}
}

func readWs(t *testing.T, conn *websocket.Conn) ws.ServerMessage {
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var serverMsg ws.ServerMessage
	err = json.Unmarshal(msg, &serverMsg)

	t.Logf("Received: %s", msg)
	return serverMsg
}

func connectWs(t *testing.T, token string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}

	return conn
}

func publishSessionStart(t *testing.T, amqpURL, sessionId string, quiz shared.Quiz) {
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	evt := shared.QuizMessage{
		SessionId: sessionId,
		Quiz:      quiz,
	}
	body, _ := json.Marshal(evt)
	ch.Publish(shared.SessionExchange, "session.start", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	rabCon.Close()
}

func publishSessionEnd(t *testing.T, amqpURL, sessionId string) {
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	body, _ := json.Marshal(sessionId)
	ch.Publish(shared.SessionExchange, "session.end", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	rabCon.Close()
}

func publishQuestionStart(t *testing.T, amqpURL, sessionId string) {
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	ch.Publish(shared.SessionExchange, fmt.Sprintf("question.%s.start", sessionId),
		false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        nil,
		})
	rabCon.Close()
}

func generateJWT(t *testing.T, session, user string, role shared.UserRole) string {
	claims := shared.UserToken{
		UserId:    user,
		UserType:  role,
		SessionId: session,
		Exp:       10000000,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Unix(0, 10000000)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rawJwt, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	require.NoError(t, err)

	return rawJwt
}

func getEnvFilePath() string {
	root, err := filepath.Abs("../../..")
	if err != nil {
		log.Fatal("failed to find project root dir")
	}
	return filepath.Join(root, ".env")
}

func startRealTimeServer(t *testing.T, amqpUrl string) {
	cfg := config.LoadConfig()

	manager := app.NewManager()

	// Connect to the rabbit MQ
	t.Log("Connecting to broker")
	err := manager.ConnectRabbitMQ(amqpUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Connected to broker")

	handlerDeps := ws.HandlerDeps{
		Tracker:  manager.QuizTracker,
		Registry: manager.ConnectionRegistry,
	}

	// SetCurrQuestionIdx route handler
	http.Handle("/ws", ws.NewWebSocketHandler(handlerDeps))

	go func() {
		err := http.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			nil)
		t.Fatal(err)
	}()

	broker, err := rabbit.NewRealTimeRabbit(manager.Rabbit)
	t.Log("Service is up!")

	go broker.ConsumeSessionStart(manager.ConnectionRegistry, manager.QuizTracker)
	go broker.ConsumeSessionEnd(manager.ConnectionRegistry, manager.QuizTracker)
	select {}
}
