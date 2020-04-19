package Worker

import (
	"go-publisher-worker/Models"
	"log"
)

type Worker struct {
	RabbitMQ *RabbitMqConnection
	env *Models.Env
}

func (worker *Worker) Work()  {
	forever := make(chan bool)

	go func() {
		for d := range worker.RabbitMQ.msgs {
			log.Printf("Received a message: %s", d.Body)




			log.Printf("Done")
			_ = d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}


func InitWorker() *Worker  {
	return &Worker{
		RabbitMQ: 				initRabbitMqConnection(Models.GetEnvStruct()),
		env:      				Models.GetEnvStruct(),
	}

}