package rabbit

// This file stores functions related to "question"-type events published to RabbitMQ

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
	"sync"
	"xxx/real_time/models"
	"xxx/real_time/ws"
	"xxx/shared"
)

// CreateQuestionStartQueue declares and binds the `question_start` queue in RabbitMQ.
// The queue is utilized to receive events "start next question".
// The queue is auto delete, since it is temporary and exists only till the session is alive.
// Returns the queue object itself, or the error if failed.
func CreateQuestionStartQueue(ch *amqp.Channel, sessionId string) (amqp.Queue, error) {
	queueName := fmt.Sprintf("question.%s.start", sessionId)
	queue, err := ch.QueueDeclare(
		queueName,
		false,
		true, // auto delete
		true,
		false,
		nil)

	if err != nil {
		return amqp.Queue{}, err
	}

	key := strings.Replace(shared.QuestionStartRoutingKey, "*", sessionId, 1)
	err = ch.QueueBind(
		queueName,
		key,
		shared.SessionExchange,
		false,
		nil)
	if err != nil {
		return amqp.Queue{}, err
	}

	return queue, nil
}

// ConsumeQuestionStart method listens to "next question start" events delivered to the corresponding queue.
func (r *RealTimeRabbit) ConsumeQuestionStart(
	registry *ws.ConnectionRegistry, tracker *ws.QuizTracker, sid string) {
	q, _ := CreateQuestionStartQueue(r.channel, sid)

	consumerTag := fmt.Sprintf("question_start_%sid", sid)

	msgs, err := r.channel.Consume(
		q.Name, // the name of the already created queue
		consumerTag,
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		fmt.Println(err)
	}

	r.QuestionStartedQsTags[sid] = consumerTag
	fmt.Println("!-----------------------! ", r.QuestionStartedQsTags)

	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Printf("Listen for new messages in question.%sid.start queue\n", sid)

	// listen to messages in parallel goroutine
	go func(s string) {
		responder := ws.NewResponder(registry, sid)

		var sessionId string

		defer wg.Done()
		for d := range msgs { // ignore the contents in the queue, since only event itself matters
			sessionId = strings.Split(d.RoutingKey, ".")[1]
			fmt.Printf("------ in consumer for %sid found sessionId %sid\n", s, sessionId)

			tracker.IncQuestionIdx(sessionId)

			qid, question := tracker.GetCurrentQuestion(sessionId)
			questionsAmount := tracker.GetQuizLen(sessionId)

			if qid == questionsAmount { // zero-based index equal to amount, means index out of range -> game already ended
				responder.SendGameEnd()
				continue
			}

			fmt.Println("next question triggered: ", qid, "in session ", sessionId)

			if qid > 0 {
				fmt.Println("Prepare Leader Board for ", sessionId)
				board, err := tracker.GetLeaderboard(sessionId)
				fmt.Println("Board from LBS: ", board)

				if err != nil {
					responder.SendError()
					fmt.Println("Leader board Error: ", err)
				} else {
					// Send LeaderBoard to Admin
					responder.SendLeaderboard(board.Table)

					allAnswers := tracker.GetAnswers(sessionId) // users' answers on all questions
					fmt.Println("USERS ANSWERS: ", allAnswers)

					currQuestionAnswers := make(map[string]models.UserAnswer)
					for userId, answers := range allAnswers {
						currQuestionAnswers[userId] = answers[qid-1]
					}

					// Send question statistics to participant
					responder.SendQuestionStat(board.Popular, currQuestionAnswers)
				}
			}

			responder.SendQuestionPayload(qid+1, // 1-based index
				questionsAmount, *question)

			// send ack for participants immediately after sending a question only for very first question
			// since further it will be sent when Admin requests it (check message.go/handleRead)
			if qid == 0 {
				responder.SendNextQuestionAck()
			}
		}
	}(sid)

	wg.Wait() // defer this function termination while consuming from the queue
	fmt.Println("Question_start queue was deleted for session ")
}

func (r *RealTimeRabbit) CleanupQuestionConsumer(sessionId string) error {
	consumerTag, ok := r.QuestionStartedQsTags[sessionId]
	if !ok {
		return fmt.Errorf("no consumer for session %s", sessionId)
	}

	err := r.channel.Cancel(consumerTag, false)
	if err != nil {
		return fmt.Errorf("failed to cancel consumer: %w", err)
	}

	delete(r.QuestionStartedQsTags, sessionId)
	return nil
}
