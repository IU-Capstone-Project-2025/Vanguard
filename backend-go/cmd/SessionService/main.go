package main

import (
	"context"
	"fmt"
	"time"
	"xxx/SessionService/Storage/Redis"
	models2 "xxx/SessionService/models"
)

func main() {
	cfg := models2.Config{
		Addr:        "localhost:6379",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}
	db, err := Redis.NewRedisClient(context.Background(), cfg)
	if err != nil {
		fmt.Println(err)
	}
	session := &models2.Session{
		ID:    "1",
		Code:  "11",
		State: "123",
	}
	err = db.SaveSession(session)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("save session success")
	}

	res, err := db.LoadSession(session.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	err = db.EditSessionState(session.ID, "idi nah")
	if err != nil {
		fmt.Println(err)
	}

	res, err = db.LoadSession(session.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}
