package Game

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"xxx/SessionService/Rabbit"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
	"xxx/SessionService/utils"
	"xxx/shared"
)

type Manager interface {
	ValidateCode(code string) bool
	GenerateUserToken(code string, UserId string, UserType shared.UserRole) *shared.UserToken
	NewSession() (*shared.Session, error)
	SessionStart(code string) error
	NextQuestion(code string) error
}

type SessionManager struct {
	rabbit     Rabbit.Broker
	cache      Redis.Cache
	codeLength int
}

func CreateSessionManager(codeLength int, rmqConn string, redisConn string) (*SessionManager, error) {
	rabbit, err := Rabbit.NewRabbit(rmqConn)
	if err != nil {
		fmt.Println("error on CreateSessionManager with rabbit", err)
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
		return &shared.Session{}, err
	}
	return session, nil
}

// ValidateCode checks that code that user sent is exist
func (manager *SessionManager) ValidateCode(code string) bool {
	flag := manager.cache.CodeExist(code)
	return flag
}

func (manager *SessionManager) GenerateUserToken(code string, UserId string, UserType shared.UserRole) *shared.UserToken {
	return &shared.UserToken{
		UserId:           UserId,
		UserType:         UserType,
		SessionId:        code,
		ServerWsEndpoint: shared.GetWsEndpoint(),
		Exp:              10000,
	}
}

func (manager *SessionManager) SessionStart(quizUUID string) error {
	//url := fmt.Sprintf("http://quiz:8001/%s", quizUUID)
	//resp, err := http.GetCurrentQuestionIdx(url)
	//if err != nil {
	//	return fmt.Errorf("error to get quiz from service %s %s", quizUUID, err.Error())
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	return fmt.Errorf("quiz session status code: %d", resp.StatusCode)
	//}
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return fmt.Errorf("error on SessionStart with read body %s, %s", quizUUID, err.Error())
	//}
	//
	var quiz shared.Quiz
	//if err := json.Unmarshal(body, &quiz); err != nil {
	//	return fmt.Errorf("error on SessionStart with unmarshal json %s %s", quizUUID, err.Error())
	//}
	err := manager.rabbit.PublishSessionStart(context.Background(), quiz)
	if err != nil {
		return fmt.Errorf("error on SessionStart with publish quiz to rabbit %s %s", quizUUID, err.Error())
	}
	return nil
}

func (manager *SessionManager) NextQuestion(code string) error {
	err := manager.rabbit.PublishQuestionStart(context.Background(), code, "aboba")
	if err != nil {
		return err
	}
	return nil
}
