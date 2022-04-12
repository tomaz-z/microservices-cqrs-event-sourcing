package mq

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/streadway/amqp"

	"eventsourcing/internal/mq/dispatcher"
	escqrs "eventsourcing/services"
)

type connector struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// Connect instantiates a new message queue connector and creates a connection, a channel and initializes queues.
func Connect(url string) (escqrs.MQConnector, error) {
	c := &connector{}

	conn, err := connection(url)
	if err != nil {
		return c, err
	}
	c.connection = conn

	ch, err := c.connection.Channel()
	if err != nil {
		return c, err
	}
	c.channel = ch

	return c, nil
}

func connection(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second * 5)
	}
	if conn == nil {
		err = errors.New("connection to mq couldn't be established")
	}
	return conn, err
}

func queues(ch *amqp.Channel, names []string) (map[string]*amqp.Queue, error) {
	queues := map[string]*amqp.Queue{}
	if names == nil {
		return queues, nil
	}

	for _, item := range names {
		q, err := ch.QueueDeclare(
			item,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return queues, err
		}
		queues[item] = &q
	}
	return queues, nil
}

func delivery(channel *amqp.Channel, queue *amqp.Queue) (<-chan amqp.Delivery, error) {
	delivery, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return delivery, err
}

func listener(handler func(event escqrs.Event), delivery <-chan amqp.Delivery) {
	go func() {
		for d := range delivery {
			event := escqrs.Event{}
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Fatal(err)
			}

			handler(event)
		}
	}()
}

func (c *connector) Close() {
	c.connection.Close()
	c.channel.Close()
}

func (c *connector) Channel() *amqp.Channel {
	return c.channel
}

func (c *connector) Dispather() escqrs.EventDispatcher {
	return dispatcher.New(c.channel)
}

func (c *connector) Start(queueNames []string, consumers []escqrs.Consumer) error {
	q, err := queues(c.channel, queueNames)
	if err != nil {
		return err
	}

	for _, consumer := range consumers {
		delivery, err := delivery(c.channel, q[consumer.Queue()])
		if err != nil {
			log.Fatal(err)
		}
		listener(consumer.Handler(), delivery)
	}

	return nil
}
