package shared

import (
	"fmt"
	"strconv"
)

type User struct {
	Name     string
	Email    string
	Password string
}

func (u User) JoinGame(code int) {
	fmt.Println(u.Name, "joined game", strconv.Itoa(code))
}
