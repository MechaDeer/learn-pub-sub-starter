package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	msg, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
	return err
}

// it does create and bind to queue, if queue already exists it will crash the client
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
	)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not subscribe to pause: %v", err)
	}

	// this needs to be in different scope, cause we want to recreate JSON
	// otherwise we might get a residual code
	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			ackType := handler(target)
			switch ackType {
			case Ack:
				msg.Ack(false)
				fmt.Println("message is acknoledged")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("message is discarded")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("message is re queued")
			}
		}
	}()

	return nil
}
