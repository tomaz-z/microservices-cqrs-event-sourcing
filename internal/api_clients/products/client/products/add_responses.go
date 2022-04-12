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

// AddReader is a Reader for the Add structure.
type AddReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewAddCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewAddBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewAddInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewAddCreated creates a AddCreated with default headers values
func NewAddCreated() *AddCreated {
	return &AddCreated{}
}

/* AddCreated describes a response with status code 201, with default header values.

Created
*/
type AddCreated struct {
	Payload *models.Product
}

func (o *AddCreated) Error() string {
	return fmt.Sprintf("[POST /][%d] addCreated  %+v", 201, o.Payload)
}
func (o *AddCreated) GetPayload() *models.Product {
	return o.Payload
}

func (o *AddCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Product)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddBadRequest creates a AddBadRequest with default headers values
func NewAddBadRequest() *AddBadRequest {
	return &AddBadRequest{}
}

/* AddBadRequest describes a response with status code 400, with default header values.

BadRequest
*/
type AddBadRequest struct {
}

func (o *AddBadRequest) Error() string {
	return fmt.Sprintf("[POST /][%d] addBadRequest ", 400)
}

func (o *AddBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewAddInternalServerError creates a AddInternalServerError with default headers values
func NewAddInternalServerError() *AddInternalServerError {
	return &AddInternalServerError{}
}

/* AddInternalServerError describes a response with status code 500, with default header values.

InternalServerError
*/
type AddInternalServerError struct {
}

func (o *AddInternalServerError) Error() string {
	return fmt.Sprintf("[POST /][%d] addInternalServerError ", 500)
}

func (o *AddInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
