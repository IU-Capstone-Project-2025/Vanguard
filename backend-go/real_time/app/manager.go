package app

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"xxx/real_time/ws"
)

// Manager represents the orchestrator of the whole service and manages the critically important components,
// as message brokers, storage, etc.
type Manager struct {
	Redis              *redis.Client
	Rabbit             *amqp.Connection
	QuizTracker        *ws.QuizTracker // map[sessionId]questionIndex
	ConnectionRegistry *ws.ConnectionRegistry
}

func NewManager() *Manager {
	return &Manager{
		Redis:              nil,
		Rabbit:             nil,
		QuizTracker:        ws.NewQuizTracker(),        // Initialize question tracker
		ConnectionRegistry: ws.NewConnectionRegistry(), // Initialize ws connections registry
	}
}

// ConnectRabbitMQ connects to the RabbitMQ using the given url
// and assigns obtained amqp.Conn to the manager.Rabbit field
func (m *Manager) ConnectRabbitMQ(url string) error {
	brokerConn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	m.Rabbit = brokerConn

	return nil
}

// ConnectRedis connects to the Redis using the given url
// and assigns obtained amqp.Conn to the manager.Redis field
func (m *Manager) ConnectRedis(url string) error {
	return nil
}
