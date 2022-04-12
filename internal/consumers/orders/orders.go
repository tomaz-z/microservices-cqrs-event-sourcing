package orders

import (
	"encoding/json"
	"log"

	escqrs "eventsourcing/services"
)

type consumer struct {
	queue       string
	eventStore  escqrs.EventStore
	ordersStore escqrs.OrdersStore
}

// New returns a new QueueOrders consumer.
func New(
	eventStore escqrs.EventStore,
	ordersStore escqrs.OrdersStore,
) escqrs.Consumer {
	return consumer{
		queue:       escqrs.QueueOrders,
		eventStore:  eventStore,
		ordersStore: ordersStore,
	}
}

func (c consumer) Queue() string {
	return c.queue
}

func (c consumer) Handler() func(event escqrs.Event) {
	return c.handler
}

func (c consumer) handler(event escqrs.Event) {
	data := escqrs.Order{
		ID: &event.AggregateID,
	}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		log.Println("error unmarshalling message: ", err)
		return
	}

	switch event.Type {
	case escqrs.EventTypeOrderAdded:
		err = c.ordersStore.AddOrder(data)
		if err != nil {
			log.Println("adding order: ", err)
			return
		}
	default:
		log.Println("received event of unsupported type")
	}

	err = c.eventStore.Apply([]escqrs.Event{event})
	if err != nil {
		log.Println("applying event: ", err)
		return
	}
}
