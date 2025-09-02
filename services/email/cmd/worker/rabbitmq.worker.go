package worker

import (
	"log"
	"time"

	"github.com/sayyidinside/monorepo-gofiber-clean/services/email/interfaces/broker/handler"
	sharedBootstrap "github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/shutdown"
)

func Consume(depedency *sharedBootstrap.Deps, shutdown *shutdown.Handler, handler *handler.Handlers) {
	q, err := depedency.RabbitMQ.Channel.QueueDeclare(
		"email", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	depedency.RabbitMQ.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	msgs, err := depedency.RabbitMQ.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			if err := handler.EmailHandler.SendEmail(d.Body); err == nil {
				log.Printf("Mail sent at: %s", time.Now())
				d.Ack(false)
			}
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	shutdown.Listen()
}
