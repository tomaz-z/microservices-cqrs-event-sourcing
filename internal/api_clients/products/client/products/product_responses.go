// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"eventsourcing/internal/api_clients/products/models"
)

// ProductReader is a Reader for the Product structure.
type ProductReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ProductReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewProductOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewProductNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewProductInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewProductOK creates a ProductOK with default headers values
func NewProductOK() *ProductOK {
	return &ProductOK{}
}

/* ProductOK describes a response with status code 200, with default header values.

OK
*/
type ProductOK struct {
	Payload *models.Product
}

func (o *ProductOK) Error() string {
	return fmt.Sprintf("[GET /{id}][%d] productOK  %+v", 200, o.Payload)
}
func (o *ProductOK) GetPayload() *models.Product {
	return o.Payload
}

func (o *ProductOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Product)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewProductNotFound creates a ProductNotFound with default headers values
func NewProductNotFound() *ProductNotFound {
	return &ProductNotFound{}
}

/* ProductNotFound describes a response with status code 404, with default header values.

NotFound
*/
type ProductNotFound struct {
}

func (o *ProductNotFound) Error() string {
	return fmt.Sprintf("[GET /{id}][%d] productNotFound ", 404)
}

func (o *ProductNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewProductInternalServerError creates a ProductInternalServerError with default headers values
func NewProductInternalServerError() *ProductInternalServerError {
	return &ProductInternalServerError{}
}

/* ProductInternalServerError describes a response with status code 500, with default header values.

InternalServerError
*/
type ProductInternalServerError struct {
}

func (o *ProductInternalServerError) Error() string {
	return fmt.Sprintf("[GET /{id}][%d] productInternalServerError ", 500)
}

func (o *ProductInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
