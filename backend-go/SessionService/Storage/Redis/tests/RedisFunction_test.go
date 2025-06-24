package tests

import (
	"context"
	"fmt"
	"testing"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
	"xxx/shared"
)

func Test_RedisFunction(t *testing.T) {
	ctx := context.Background()
	RedisConfig := models.Config{
		Addr:        "localhost:6379",
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
