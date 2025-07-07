package rabbit

// This file stores functions related to "question"-type events published to RabbitMQ

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
	"sync"
	"xxx/real_time/ws"
	"xxx/shared"
)

// CreateQuestionStartQueue declares and binds the `question_start` queue in RabbitMQ.
// The queue is utilized to receive events "start next question".
// The queue is auto delete, since it is temporary and exists only till the session is alive.
// Returns the queue object itself, or the error if failed.
func CreateQuestionStartQueue(ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		"",
		false,
		true, // auto delete
		true,
		false,
		nil)

	if err != nil {
		return amqp.Queue{}, err
	}

	err = ch.QueueBind(
		queue.Name,
		shared.QuestionStartRoutingKey,
		shared.SessionExchange,
		false,
		nil)
	if err != nil {
		return amqp.Queue{}, err
	}

	return queue, nil
}

// ConsumeQuestionStart method listens to "next question start" events delivered to the corresponding queue.
func (r *RealTimeRabbit) ConsumeQuestionStart(registry *ws.ConnectionRegistry, tracker *ws.QuizTracker) {
	q, _ := CreateQuestionStartQueue(r.channel)

	msgs, err := r.channel.Consume(
		q.Name, // the name of the already created queue
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Printf("Listen for new messages in question.*.start queue\n")

	// listen to messages in parallel goroutine
	go func() {
		var sessionId string

		defer wg.Done()
		for d := range msgs { // ignore the contents in the queue, since only event itself matters
			sessionId = strings.Split(d.RoutingKey, ".")[1]

			tracker.IncQuestionIdx(sessionId)

			qid, question := tracker.GetCurrentQuestion(sessionId)

			fmt.Println("next question triggered: ", qid, "in session ", sessionId)

			questionPayloadMsg := ws.ServerMessage{
				Type:            ws.MessageTypeQuestion,
				QuestionIdx:     qid + 1,
				QuestionsAmount: tracker.GetQuizLen(sessionId),
				Text:            question.Text,
				Options:         question.Options,
			}

			registry.SendToAdmin(sessionId, questionPayloadMsg.Bytes())

			nextQuestionAck := ws.ServerMessage{
					Type:          ws.MessageTypeNextQuestion,
				}
			registry.BroadcastToSession(sessionId, nextQuestionAck.Bytes(), false)

			if qid == 0 {
				gameStartAck := ws.ServerMessage{
					Type:          ws.MessageTypeAck,
					IsGameStarted: true,
				}
				registry.BroadcastToSession(sessionId, gameStartAck.Bytes(), false)
			}
		}
		
	}()

	wg.Wait() // defer this function termination while consuming from the queue
}
