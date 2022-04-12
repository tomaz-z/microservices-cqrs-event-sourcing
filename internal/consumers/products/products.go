package products

import (
	"encoding/json"
	"log"

	escqrs "eventsourcing/services"
)

type consumer struct {
	queue         string
	eventStore    escqrs.EventStore
	productsStore escqrs.ProductsStore
}

// New returns a new QueueProducts consumer.
func New(
	eventStore escqrs.EventStore,
	productsStore escqrs.ProductsStore,
) escqrs.Consumer {
	return consumer{
		queue:         escqrs.QueueProducts,
		eventStore:    eventStore,
		productsStore: productsStore,
	}
}

func (c consumer) Queue() string {
	return c.queue
}

func (c consumer) Handler() func(event escqrs.Event) {
	return c.handler
}

func (c consumer) handler(event escqrs.Event) {
	data := escqrs.Product{
		ID: &event.AggregateID,
	}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		log.Println("error unmarshalling message: ", err)
		return
	}

	switch event.Type {
	case escqrs.EventTypeProductAdded:
		err := c.productsStore.AddProduct(data)
		if err != nil {
			log.Println("adding product: ", err)
			return
		}
	case escqrs.EventTypeProductQuantityUpdated:
		err := c.productsStore.UpdateProductQuantity(data)
		if err != nil {
			log.Println("updating product quantity: ", err)
			return
		}
	default:
		log.Println("received event of unsupported type")
	}

	err = c.eventStore.Apply([]escqrs.Event{event})
	if err != nil {
		log.Println("applying event: ", err)
	}
}
