package dispatcher

import (
	"encoding/json"

	"github.com/streadway/amqp"

	escqrs "eventsourcing/services"
)

type dispatcher struct {
	channel *amqp.Channel
}

// New returns a new instance of an event dispatcher.
func New(
	channel *amqp.Channel,
) escqrs.EventDispatcher {
	return &dispatcher{
		channel: channel,
	}
}

func (d *dispatcher) Dispatch(
	queueName string,
	event escqrs.Event,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return d.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
