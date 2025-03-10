package pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func Subscribe[T any](
	conn *amqp.Connection,
	exchange,
	exchangeType,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, exchangeType, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("couldn't declare and bind queue: %v", err)
	}

	err = ch.Qos(5, 0, false)
	if err != nil {
		return fmt.Errorf("could not set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			var dat T
			err = json.Unmarshal(msg.Body, &dat)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			whatack := handler(dat)
			switch whatack {
			case Ack:
				if err = msg.Ack(false); err != nil {
					log.Printf("couldn't Ack the message: %v", err)
					continue
				}
			case NackRequeue:
				if err = msg.Nack(false, true); err != nil {
					log.Printf("couldn't Nack Requeue the message: %v", err)
					continue
				}
			case NackDiscard:
				if err = msg.Nack(false, false); err != nil {
					log.Printf("couldn't Nack Discard the message: %v", err)
					continue
				}
			}
		}
	}()
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	exchangeType,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("couldn't create channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		exchange,
		exchangeType,
		simpleQueueType == SimpleQueueDurable,
		simpleQueueType == SimpleQueueTransient,
		false,
		false,
		nil)
	if err != nil {
		// Check if error is due to conflict, otherwise return error
		if strings.Contains(err.Error(), "inequivalent arg") {
			// Handle property conflict case
			return nil, amqp.Queue{}, fmt.Errorf("exchange exists with different properties: %v", err)
		}
		// Normal creation error
		return nil, amqp.Queue{}, fmt.Errorf("couldn't declare exchange: %v", err)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		simpleQueueType == SimpleQueueDurable,
		simpleQueueType == SimpleQueueTransient,
		simpleQueueType == SimpleQueueTransient,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "crawler_dlx",
		},
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("couldn't declare queue: %v", err)
	}

	err = ch.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("couldn't bind queue: %v", err)
	}

	return ch, queue, nil
}

func DeclareDeadLetterSetUp(connc *amqp.Connection) error {
	ch, err := connc.Channel()
	if err != nil {
		return fmt.Errorf("couldn't open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"crawler_dlx",
		"fanout",
		true,  //durable
		false, //auto-delete
		false, //internal
		false, //no-wait
		nil,   //arguments
	)
	if err != nil {
		return fmt.Errorf("couldn't declare dead letter exchange: %v", err)
	}

	dlq, err := ch.QueueDeclare(
		"crawler_dlq",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("couldn't declare dead letter queue: %v", err)
	}

	err = ch.QueueBind(
		dlq.Name,
		"", //fanout no key
		"crawler_dlx",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("couldn't bind dead letter queue: %v", err)
	}
	return nil

}
