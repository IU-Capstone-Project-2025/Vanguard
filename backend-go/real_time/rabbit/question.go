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
	registry *ws.ConnectionRegistry, tracker *ws.QuizTracker, s string) {
	q, _ := CreateQuestionStartQueue(r.channel, s)

	consumerTag := fmt.Sprintf("question_start_%s", s)

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

	r.QuestionStartedQsTags[s] = consumerTag
	fmt.Println("!-----------------------! ", r.QuestionStartedQsTags)

	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Printf("Listen for new messages in question.%s.start queue\n", s)

	// listen to messages in parallel goroutine
	go func(s string) {
		var sessionId string

		defer wg.Done()
		for d := range msgs { // ignore the contents in the queue, since only event itself matters
			sessionId = strings.Split(d.RoutingKey, ".")[1]
			fmt.Printf("------ in consumer for %s found sessionId %s\n", s, sessionId)

			tracker.IncQuestionIdx(sessionId)

			qid, question := tracker.GetCurrentQuestion(sessionId)
			questionsAmount := tracker.GetQuizLen(sessionId)

			if qid == questionsAmount { // zero-based index equal to amount, means index out of range -> game already ended
				gameEndAck := ws.ServerMessage{
					Type: ws.MessageTypeEnd,
				}
				registry.BroadcastToSession(sessionId, gameEndAck.Bytes(), false)
				continue
			}

			fmt.Println("next question triggered: ", qid, "in session ", sessionId)

			if qid > 0 {
				fmt.Println("Prepare Leader Board for ", sessionId)
				board, err := tracker.GetLeaderboard(sessionId)
				fmt.Println("Board from LBS: ", board)

				if err != nil {
					errorResp := ws.ServerMessage{Type: ws.MessageTypeError}
					registry.BroadcastToSession(sessionId, errorResp.Bytes(), true)
					fmt.Println("Leader board Error: ", err)
				} else {
					// Send LeaderBoard to Admin
					leaderBoard := ws.ServerMessage{
						Type:    ws.MessageTypeLeaderboard,
						Payload: board.Table,
					}

					registry.SendToAdmin(sessionId, leaderBoard.Bytes())

					usersAnswers := tracker.GetAnswers(sessionId)
					fmt.Println("USERS ANSWERS: ", usersAnswers)

					// Send question statistics to participant
					for _, connectionCtx := range registry.GetConnections(sessionId) {
						if connectionCtx.Role == shared.RoleAdmin {
							continue
						}
						user := connectionCtx.UserId

						stat := ws.ServerMessage{
							Type:    ws.MessageTypeStat,
							Correct: usersAnswers[user][qid-1].Correct,
							Payload: board.Popular,
						}

						registry.SendMessage(stat.Bytes(), connectionCtx)

					}
				}

			}

			questionPayloadMsg := ws.ServerMessage{
				Type:            ws.MessageTypeQuestion,
				QuestionIdx:     qid + 1,
				QuestionsAmount: questionsAmount,
				Text:            question.Text,
				Options:         question.Options,
			}

			registry.SendToAdmin(sessionId, questionPayloadMsg.Bytes())

			fmt.Printf("Send question payload for %s: %v\n", sessionId, questionPayloadMsg)

			nextQuestionAck := ws.ServerMessage{
				Type: ws.MessageTypeNextQuestion,
			}
			registry.BroadcastToSession(sessionId, nextQuestionAck.Bytes(), false)

			fmt.Printf("Send next question ack for %s: %v\n", sessionId, nextQuestionAck)
		}
	}(s)

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
