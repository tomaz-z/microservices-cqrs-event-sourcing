package aggregates

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	escqrs "eventsourcing/services"
)

type aggregate struct {
	id     *uuid.UUID
	domain string
	data   interface{}
}

// New returns a new aggregate entity.
func New(domain string) escqrs.Aggregate {
	return &aggregate{
		domain: domain,
	}
}

func (a *aggregate) WithID(id uuid.UUID) escqrs.Aggregate {
	a.id = &id
	return a
}

func (a *aggregate) ID() *uuid.UUID {
	return a.id
}

func (a *aggregate) Handle(command escqrs.Command) (escqrs.Event, error) {
	if a.domain != command.Domain() {
		return escqrs.Event{}, errors.New("domain name of the aggregate and the command do not match")
	}

	if a.id == nil {
		id := uuid.New()
		a.id = &id
	}

	event := escqrs.Event{
		ID:            uuid.New(),
		AggregateID:   *a.id,
		AggregateType: a.domain,
		CreatedAt:     time.Now(),
	}

	var err error
	switch command.Type() {
	case escqrs.CommandTypeAddProduct:
		event.Type = escqrs.EventTypeProductAdded
	case escqrs.CommandTypeUpdateProductQuantity:
		event.Type = escqrs.EventTypeProductQuantityUpdated
	case escqrs.CommandTypeAddOrder:
		event.Type = escqrs.EventTypeOrderAdded
	case escqrs.CommandTypeReserveProducts:
		event.Type = escqrs.EventTypeProductsReserved
	}

	event.Data, err = json.Marshal(command.Data())
	if err != nil {
		return event, err
	}

	a.data = command.Data()

	return event, nil
}
