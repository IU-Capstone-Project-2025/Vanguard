package Game

import (
	"github.com/google/uuid"
	"xxx/SessionService/Rabbit"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
	"xxx/SessionService/utils"
)

type Manager interface {
	ValidateCode(code string) bool
	GenerateToken(code string) *models.UserToken
	NewSession(codeLength int) (*models.Session, error)
	PublishEvent(payload interface{}) error
}

type SessionManager struct {
	rabbit     Rabbit.Rabbit
	cache      Redis.Cache
	codeLength int
}

func CreateSessionManager(cache Redis.Cache, codeLength int) *SessionManager {
	return &SessionManager{
		cache:      cache,
		codeLength: codeLength,
	}
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
		ID:    sessionId,
		Code:  code,
		State: "waiting",
	}
	err := manager.cache.SaveSession(session)
	if err != nil {
		return &models.Session{}, err
	}
	err = manager.rabbit.SessionCreated(session)
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

func (manager *SessionManager) GenerateToken(code string, UserId string) *models.UserToken {
	return &models.UserToken{
		UserId:      UserId,
		CurrentQuiz: code,
		Exp:         10000,
	}
}

func (manager *SessionManager) PublishEvent(payload interface{}) error {
	err := manager.rabbit.PublishEvent(payload)
	if err != nil {
		return err
	}
	return nil
}
