package reserveproducts

import (
	"encoding/json"
	"log"

	"eventsourcing/internal/aggregates"
	"eventsourcing/internal/commands"
	escqrs "eventsourcing/services"
)

type consumer struct {
	queue         string
	eventStore    escqrs.EventStore
	productsStore escqrs.ProductsStore
	dispatcher    escqrs.EventDispatcher
}

// New returns a new QueueReserveProducts consumer.
func New(
	eventStore escqrs.EventStore,
	productsStore escqrs.ProductsStore,
	dispatcher escqrs.EventDispatcher,
) escqrs.Consumer {
	return consumer{
		queue:         escqrs.QueueReserveProducts,
		eventStore:    eventStore,
		productsStore: productsStore,
		dispatcher:    dispatcher,
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
		log.Println("unmarshalling event: ", err)
		return
	}

	switch event.Type {
	case escqrs.EventTypeProductsReserved:
		products, err := c.productsStore.ReserveProducts(data.Products)
		if err != nil {
			log.Println("reserving products :", err)
			return
		}
		for i, item := range products {
			data.Products[i].Price = &item.Price
		}

		eventAddOrder, err := aggregates.New(
			escqrs.DomainOrder,
		).Handle(
			commands.New(
				escqrs.DomainOrder,
				escqrs.CommandTypeAddOrder,
				data,
			),
		)
		if err != nil {
			log.Println("creating event: ", err)
			return
		}

		err = c.dispatcher.Dispatch(
			escqrs.QueueOrders,
			eventAddOrder,
		)
		if err != nil {
			log.Println("dispatching message to orders", err)
			return
		}

		err = c.eventStore.Apply([]escqrs.Event{event})
		if err != nil {
			log.Println("applying event: ", err)
			return
		}
	default:
		log.Println("received event of unsupported type")
	}
}
