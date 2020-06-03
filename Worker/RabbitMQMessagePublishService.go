package Worker

import (
	"fmt"
	"github.com/streadway/amqp"
	"go-publisher-worker/Models"
)

type RabbitMQPublisher struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	q    amqp.Queue
}

func (rabbitMQ *RabbitMQPublisher) publishMessageOnQueue(message []byte) error {
	err := rabbitMQ.Ch.Publish(
		"",
		rabbitMQ.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			Body:            message,
		})
	if err != nil {
		return err
	}
	return nil
}

func InitRabbitMQVideoAnalysisPublisher() *RabbitMQPublisher  {
	fmt.Println("Video publiseh queue start connection")
	env := Models.GetEnvStruct()
	conn, err := amqp.Dial("amqp://" + env.RabbitUser + ":" + env.RabbitPassword + "@" + env.RabbitHost + ":" + env.RabbitPort)
	failOnError(err, "Failed to connect to RabbitMQ video publisher")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel video publisher")

	q, err := ch.QueueDeclare(
		env.RabbitMqVideoAnalysis, // name
		true,         	// durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue video publisher")
	fmt.Println("Video publiseh queue end")
	return &RabbitMQPublisher{
		Conn: conn,
		Ch:   ch,
		q:    q,
	}
}