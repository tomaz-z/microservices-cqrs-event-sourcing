package patch

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"

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

// New returns a new patch product handler.
func New(
	eventStore escqrs.EventStore,
	dispatcher escqrs.EventDispatcher,
) productsOperations.PatchHandler {
	return handler{
		eventStore: eventStore,
		dispatcher: dispatcher,
	}
}

func (h handler) Handle(params productsOperations.PatchParams) middleware.Responder {
	product := escqrs.ProductPatch{
		ID:       uuid.Must(uuid.Parse(params.ID.String())),
		Quantity: params.Product.Quantity,
	}

	event, err := aggregates.New(
		escqrs.DomainProduct,
	).WithID(uuid.MustParse(params.ID.String())).
		Handle(commands.New(
			escqrs.DomainProduct,
			escqrs.CommandTypeUpdateProductQuantity,
			product,
		))
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

	return productsOperations.NewPatchOK().WithPayload(&models.PatchProduct{
		Quantity: product.Quantity,
	})
}
