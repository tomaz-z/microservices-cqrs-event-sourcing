// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ProductHandlerFunc turns a function with the right signature into a product handler
type ProductHandlerFunc func(ProductParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ProductHandlerFunc) Handle(params ProductParams) middleware.Responder {
	return fn(params)
}

// ProductHandler interface for that can handle valid product params
type ProductHandler interface {
	Handle(ProductParams) middleware.Responder
}

// NewProduct creates a new http.Handler for the product operation
func NewProduct(ctx *middleware.Context, handler ProductHandler) *Product {
	return &Product{Context: ctx, Handler: handler}
}

/* Product swagger:route GET /{id} products product

Get product by ID.

*/
type Product struct {
	Context *middleware.Context
	Handler ProductHandler
}

func (o *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewProductParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
