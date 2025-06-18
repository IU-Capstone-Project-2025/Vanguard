package Game

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"xxx/SessionService/Rabbit"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
	"xxx/SessionService/utils"
	"xxx/Shared"
)

type Manager interface {
	ValidateCode(code string) bool
	GenerateUserToken(code string, UserId string, UserType string) *models.UserToken
	NewSession() (*models.Session, error)
	PublishEvent(payload interface{}) error
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
func (manager *SessionManager) NewSession() (*models.Session, error) {
	sessionId := uuid.New().String()
	code := ""
	for i := 0; i < 3; i++ {
		code = utils.GenerateSessionCode(manager.codeLength)
		if !manager.cache.CodeExist(code) {
			break
		}
	}
	//TODO improve session code generation
	session := &models.Session{
		ID:               sessionId,
		Code:             code,
		State:            "waiting",
		ServerWsEndpoint: Shared.WsEndpoint,
	}
	err := manager.cache.SaveSession(session)
	if err != nil {
		return &models.Session{}, err
	}
	err = manager.rabbit.PublishEvent(session)
	if err != nil {
		return &models.Session{}, err
	}
	return session, nil
}

// ValidateCode checks that code that user sent is exist
func (manager *SessionManager) ValidateCode(code string) bool {
	flag := manager.cache.CodeExist(code)
	return flag
}

func (manager *SessionManager) GenerateUserToken(code string, UserId string, UserType string) *models.UserToken {
	return &models.UserToken{
		UserId:           UserId,
		UserType:         UserType,
		CurrentQuiz:      code,
		ServerWsEndpoint: Shared.WsEndpoint,
		Exp:              10000,
	}
}

func (manager *SessionManager) PublishEvent(payload interface{}) error {
	err := manager.rabbit.PublishEvent(payload)
	if err != nil {
		return err
	}
	return nil
}
