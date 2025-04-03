package pubsub

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable = iota
	SimpleQueueTransient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	ch.Qos(10, 0, false)

	var durable bool
	var autoDelete bool
	var exclusive bool
	switch simpleQueueType {
	case SimpleQueueDurable:
		durable = true
		autoDelete = false
		exclusive = false
	case SimpleQueueTransient:
		durable = false
		autoDelete = true
		exclusive = true
	default:
		return nil, amqp.Queue{}, errors.New("invalid queue type")
	}

	table := amqp.Table{}
	// value for it is the name of the exchange we want to use as dead letter
	table["x-dead-letter-exchange"] = "peril_dlx"
	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, table)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(
		queue.Name, // queue name
		key,        // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return ch, queue, nil
}
