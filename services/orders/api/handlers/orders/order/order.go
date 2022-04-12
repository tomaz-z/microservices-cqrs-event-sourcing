package order

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	escqrs "eventsourcing/services"
	"eventsourcing/services/orders/api/models"
	"eventsourcing/services/orders/api/restapi/operations/orders"
)

type handler struct {
	ordersStore escqrs.OrdersStore
}

// New returns a new handler for GET /orders/{id}.
func New(ordersStore escqrs.OrdersStore) orders.OrderHandler {
	return handler{
		ordersStore: ordersStore,
	}
}

func (h handler) Handle(params orders.OrderParams) middleware.Responder {
	order, err := h.ordersStore.GetOrder(uuid.MustParse(params.ID.String()))
	if err != nil {
		log.Println(err)
		return orders.NewOrderInternalServerError()
	}

	products := []*models.Product{}
	for _, item := range order.Products {
		products = append(products, &models.Product{
			ID:       strfmt.UUID(item.ID.String()),
			Quantity: item.Quantity,
		})
	}

	return orders.NewOrderOK().WithPayload(&models.Order{
		ID:       strfmt.UUID(order.ID.String()),
		Products: products,
	})
}
