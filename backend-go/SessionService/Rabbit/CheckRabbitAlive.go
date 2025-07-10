package Rabbit

import (
	"fmt"
)

func (r *Rabbit) CheckRabbitAlive() error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbit error channel %v", err)
	}
	defer ch.Close()
	return nil
}
