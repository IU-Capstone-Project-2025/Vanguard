package Game

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"xxx/SessionService/Rabbit"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
	"xxx/SessionService/utils"
	"xxx/shared"
)

type Manager interface {
	ValidateCode(code string) bool
	GenerateUserToken(code string, UserId string, UserName string, UserType shared.UserRole) *shared.UserToken
	NewSession() (*shared.Session, error)
	SessionStart(quizUUID string, sessionId string) error
	NextQuestion(code string) error
	GetListOfUsers(quizUUID string) ([]string, error)
	AddPlayerToSession(quizUUID string, UserName string) error
	SessionStartMock(quizUUID string, sessionId string) error
}

type SessionManager struct {
	rabbit     Rabbit.Broker
	cache      Redis.Cache
	codeLength int
}

func CreateSessionManager(codeLength int, rmqConn string, redisConn string) (*SessionManager, error) {
	rabbit, err := Rabbit.NewRabbit(rmqConn)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	RedisConfig := models.Config{
		Addr:        redisConn,
		Password:    "",
		DB:          0,
		MaxRetries:  0,
		DialTimeout: 0,
		Timeout:     0,
	}
	redis, err := Redis.NewRedisClient(ctx, RedisConfig)
	if err != nil {
		fmt.Println("error on CreateSessionManager with redis", err)
		return nil, err
	}
	fmt.Println("Create Session manager ok")
	return &SessionManager{
		rabbit:     rabbit,
		cache:      redis,
		codeLength: codeLength,
	}, nil
}

// NewSession create session and save it to Redis
func (manager *SessionManager) NewSession() (*shared.Session, error) {
	sessionId := uuid.New().String()
	code := ""
	for i := 0; i < 3; i++ {
		code = utils.GenerateSessionCode(manager.codeLength)
		if !manager.cache.CodeExist(code) {
			break
		}
	}
	//TODO improve session code generation
	session := &shared.Session{
		ID:               sessionId,
		Code:             code,
		State:            "waiting",
		ServerWsEndpoint: shared.GetWsEndpoint(),
	}
	err := manager.cache.SaveSession(session)
	if err != nil {
		return &shared.Session{}, fmt.Errorf("error saving session to redis: %v", err)
	}
	return session, nil
}

// ValidateCode checks that code that user sent is exist
func (manager *SessionManager) ValidateCode(code string) bool {
	flag := manager.cache.CodeExist(code)
	return flag
}

func (manager *SessionManager) GenerateUserToken(code string, UserId string, UserName string, UserType shared.UserRole) *shared.UserToken {
	return &shared.UserToken{
		UserId:    UserId,
		UserType:  UserType,
		SessionId: code,
		Exp:       10000,
	}
}

func (manager *SessionManager) SessionStart(quizUUID string, sessionId string) error {
	fmt.Println(quizUUID)
	url := fmt.Sprintf("%s%s", shared.QuizManager, quizUUID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error to get quiz from service %s %s", quizUUID, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("quiz session status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error on SessionStart with read body %s, %s", quizUUID, err.Error())
	}

	var quiz shared.Quiz
	if err := json.Unmarshal(body, &quiz); err != nil {
		return fmt.Errorf("error on SessionStart with unmarshal json %s %s", quizUUID, err.Error())
	}
	message := shared.QuizMessage{
		SessionId: sessionId,
		Quiz:      quiz,
	}
	err = manager.rabbit.PublishSessionStart(context.Background(), message)
	if err != nil {
		return fmt.Errorf("error on SessionStart with publish quiz to rabbit %s %s", quizUUID, err.Error())
	}
	return nil
}

func (manager *SessionManager) NextQuestion(code string) error {
	err := manager.rabbit.PublishQuestionStart(context.Background(), code, "aboba")
	if err != nil {
		return fmt.Errorf("error to send message to rabbit %s", err)
	}
	return nil
}
func (manager *SessionManager) AddPlayerToSession(quizUUID string, UserName string) error {
	err := manager.cache.AddPlayerToSession(quizUUID, UserName)
	if err != nil {
		return fmt.Errorf("error saving user to redis: %v", err)
	}
	return nil
}

func (manager *SessionManager) GetListOfUsers(quizUUID string) ([]string, error) {
	users, err := manager.cache.GetPlayersForSession(quizUUID)
	if err != nil {
		return nil, fmt.Errorf("error get user from redis: %v", err)
	}
	return users, nil
}

func (manager *SessionManager) SessionStartMock(quizUUID string, sessionId string) error {
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
	message := shared.QuizMessage{
		SessionId: sessionId,
		Quiz:      quiz,
	}
	err := manager.rabbit.PublishSessionStart(context.Background(), message)
	if err != nil {
		return fmt.Errorf("error on SessionStart with publish quiz to rabbit %s %s", quizUUID, err.Error())
	}
	return nil
}
