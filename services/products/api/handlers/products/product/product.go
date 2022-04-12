package product

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	escqrs "eventsourcing/services"
	"eventsourcing/services/products/api/models"
	productsOps "eventsourcing/services/products/api/restapi/operations/products"
)

type handler struct {
	productsStore escqrs.ProductsStore
}

// New returns a new handler for GET /products/{id}.
func New(productsStore escqrs.ProductsStore) productsOps.ProductHandler {
	return handler{
		productsStore: productsStore,
	}
}

func (h handler) Handle(params productsOps.ProductParams) middleware.Responder {
	product, err := h.productsStore.GetProduct(uuid.Must(uuid.Parse(params.ID.String())))
	if err != nil {
		if err.Error() == "no item found" {
			return productsOps.NewProductNotFound()
		}

		log.Println(err)
		return productsOps.NewProductInternalServerError()
	}
	return productsOps.NewProductOK().WithPayload(&models.Product{
		ID:          strfmt.UUID(product.ID.String()),
		Name:        &product.Name,
		Description: &product.Description,
		Price:       &product.Price,
		Quantity:    &product.Quantity,
	})
}
