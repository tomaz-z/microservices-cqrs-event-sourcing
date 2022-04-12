package products

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	escqrs "eventsourcing/services"
	"eventsourcing/services/products/api/models"
	productsOps "eventsourcing/services/products/api/restapi/operations/products"
)

type handler struct {
	productsStore escqrs.ProductsStore
}

// New returns a new handler for GET /products
func New(productsStore escqrs.ProductsStore) productsOps.ProductsHandler {
	return handler{
		productsStore: productsStore,
	}
}

func (h handler) Handle(params productsOps.ProductsParams) middleware.Responder {
	products, err := h.productsStore.GetProducts()
	if err != nil {
		return productsOps.NewProductsInternalServerError()
	}

	return productsOps.NewProductsOK().
		WithPayload(
			h.productsToModels(products),
		)
}

func (h handler) productsToModels(products []escqrs.Product) models.ProductList {
	productList := models.ProductList{}
	for _, item := range products {
		productList = append(productList, &models.Product{
			ID:          strfmt.UUID(item.ID.String()),
			Name:        &item.Name,
			Description: &item.Description,
			Price:       &item.Price,
			Quantity:    &item.Quantity,
		})
	}
	return productList
}
