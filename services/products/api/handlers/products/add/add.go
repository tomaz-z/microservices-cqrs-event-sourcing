package add

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"eventsourcing/internal/aggregates"
	"eventsourcing/internal/commands"
	escqrs "eventsourcing/services"
	"eventsourcing/services/products/api/models"
	productsOperations "eventsourcing/services/products/api/restapi/operations/products"
)

type handler struct {
	eventStore escqrs.EventStore
	dispatcher escqrs.EventDispatcher
}

// New returns new POST /products handler.
func New(
	eventStore escqrs.EventStore,
	dispatcher escqrs.EventDispatcher,
) productsOperations.AddHandler {
	return handler{
		eventStore: eventStore,
		dispatcher: dispatcher,
	}
}

func (h handler) Handle(params productsOperations.AddParams) middleware.Responder {
	product := escqrs.Product{
		Name:        *params.Product.Name,
		Description: *params.Product.Description,
		Quantity:    *params.Product.Quantity,
		Price:       *params.Product.Price,
	}

	event, err := aggregates.New(
		escqrs.DomainProduct,
	).Handle(
		commands.New(
			escqrs.DomainProduct,
			escqrs.CommandTypeAddProduct,
			product,
		),
	)
	if err != nil {
		log.Println(err)
		return productsOperations.NewAddInternalServerError()
	}

	err = h.dispatcher.Dispatch(
		escqrs.QueueProducts,
		event,
	)
	if err != nil {
		log.Println(err)
		return productsOperations.NewAddInternalServerError()
	}

	return productsOperations.NewAddCreated().WithPayload(
		&models.Product{
			ID:          strfmt.UUID(event.AggregateID.String()),
			Name:        &product.Name,
			Description: &product.Description,
			Quantity:    &product.Quantity,
			Price:       &product.Price,
		},
	)
}
