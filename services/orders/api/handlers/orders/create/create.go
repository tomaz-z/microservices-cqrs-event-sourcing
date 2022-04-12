package create

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	"eventsourcing/internal/aggregates"
	"eventsourcing/internal/commands"
	escqrs "eventsourcing/services"
	"eventsourcing/services/orders/api/models"
	ordersOperations "eventsourcing/services/orders/api/restapi/operations/orders"
)

type handler struct {
	dispatcher escqrs.EventDispatcher
}

// New returns a new POST /products handler.
func New(
	dispatcher escqrs.EventDispatcher,
) ordersOperations.CreateHandler {
	return handler{
		dispatcher: dispatcher,
	}
}

func (h handler) Handle(params ordersOperations.CreateParams) middleware.Responder {
	products := []escqrs.OrderProduct{}
	for _, item := range params.Order.Products {
		id := uuid.Must(uuid.Parse(item.ID.String()))

		products = append(products, escqrs.OrderProduct{
			ID:       id,
			Quantity: item.Quantity,
		})
	}

	event, err := aggregates.New(
		escqrs.DomainProduct,
	).Handle(
		commands.New(
			escqrs.DomainProduct,
			escqrs.CommandTypeReserveProducts,
			escqrs.Order{
				Products: products,
			},
		),
	)
	if err != nil {
		log.Println(err)
		return ordersOperations.NewCreateInternalServerError()
	}

	err = h.dispatcher.Dispatch(
		escqrs.QueueReserveProducts,
		event,
	)
	if err != nil {
		log.Println(err)
		return ordersOperations.NewCreateInternalServerError()
	}

	return ordersOperations.NewCreateCreated().WithPayload(
		&models.Order{
			ID:       strfmt.UUID(event.AggregateID.String()),
			Products: params.Order.Products,
		},
	)
}
