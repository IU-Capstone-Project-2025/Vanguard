package app

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/url"
	"strconv"
	"xxx/real_time/cache"
	"xxx/real_time/cache/redis"
	"xxx/real_time/rabbit"
	"xxx/real_time/ws"
)

// Manager represents the orchestrator of the whole service and manages the critically important components,
// as message brokers, storage, etc.
type Manager struct {
	Redis              cache.Cache
	Rabbit             *rabbit.RealTimeRabbit
	QuizTracker        *ws.QuizTracker // map[sessionId]questionIndex
	ConnectionRegistry *ws.ConnectionRegistry
}

func NewManager(lbHost, lbPort string) *Manager {
	leaderboardUrl := fmt.Sprintf("%s:%s", lbHost, lbPort)

	return &Manager{
		Redis:              nil,
		Rabbit:             nil,
		QuizTracker:        ws.NewQuizTracker(leaderboardUrl), // Initialize question tracker
		ConnectionRegistry: ws.NewConnectionRegistry(),        // Initialize ws connections registry
	}
}

// ConnectRabbitMQ connects to the RabbitMQ using the given url
// and assigns obtained amqp.Conn to the manager.Rabbit field
func (m *Manager) ConnectRabbitMQ(url string) (*rabbit.RealTimeRabbit, error) {
	brokerConn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	broker, err := rabbit.NewRealTimeRabbit(brokerConn)

	m.Rabbit = broker

	return broker, nil
}

// ConnectRedis connects to the Redis using the given url
// and assigns obtained amqp.Conn to the manager.Redis field
func (m *Manager) ConnectRedis(rawUrl string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	addr := u.Host               // host:port
	pass, _ := u.User.Password() // may be empty
	db := 0                      // default DB
	if len(u.Path) > 1 {         // "/2" in path -> DB = 2
		if n, err := strconv.Atoi(u.Path[1:]); err == nil {
			db = n
		}
	}

	client := redis.NewClient(addr, pass, db)
	if err := client.Ping(); err != nil { // fail fast
		return err
	}

	m.Redis = client
	m.QuizTracker.SetCache(client)

	return nil
}
