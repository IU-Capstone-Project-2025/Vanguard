package tests

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"testing"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
	"xxx/shared"
)

func getEnvFilePath() string {
	envPath := filepath.Join("..", "..", "..", "..", "..", ".env") // сдвигаемся на 4 уровня вверх из integration_tests
	absPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}
func Test_RedisFunction(t *testing.T) {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	redisUrl := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	ctx := context.Background()
	RedisConfig := models.Config{
		Addr:        redisUrl,
		Password:    "",
		DB:          0,
		MaxRetries:  0,
		DialTimeout: 0,
		Timeout:     0,
	}
	redis, err := Redis.NewRedisClient(ctx, RedisConfig)
	if err != nil {
		t.Error("error creating redis client", "error", err)
	}
	err = redis.SaveSession(&shared.Session{
		ID:               "testRedis",
		Code:             "testRedis",
		State:            "testRedis",
		ServerWsEndpoint: "testRedis",
	})
	if err != nil {
		t.Error("error creating redis client", "error", err)
	}
	fmt.Println("SaveSession success")
	session, err := redis.LoadSession("testRedis")
	if err != nil {
		t.Error("error loading session", "error", err)
	}
	if session.Code != "testRedis" {
		t.Error("error loading session", "error", err)
	}
	fmt.Println("LoadSession success", session)
	err = redis.EditSessionState("testRedis", "testRedis2")
	if err != nil {
		t.Error("error editing session", "error", err)
	}
	session, err = redis.LoadSession("testRedis")
	if err != nil {
		t.Error("error loading session", "error", err)
	}
	if session.State != "testRedis2" {
		t.Error("error loading session", "error", err)
	}
	fmt.Println("EditSessionState success", session)
	err = redis.DeleteSession("testRedis")
	if err != nil {
		t.Error("error deleting session", "error", err)
	}
	session, err = redis.LoadSession("testRedis")
	if err != nil {
		fmt.Println("DeleteSession success")
	}

}
