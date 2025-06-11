package shared

import (
	"fmt"
	"math/rand"
)

type Admin struct {
	Name     string
	Email    string
	Password string
}

func (a Admin) CreateGame() int {
	minn := 100000
	maxx := 999999
	n := rand.Intn(maxx-minn+1) + minn
	fmt.Println(a.Name, "created game", n)
	return n

}
