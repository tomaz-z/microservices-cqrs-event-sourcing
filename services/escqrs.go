package escqrs

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Events.
const (
	// Products.
	EventTypeProductAdded           = "product-added"
	EventTypeProductQuantityUpdated = "product-quantity-updated"

	// Orders.
	EventTypeOrderAdded       = "order-added"
	EventTypeProductsReserved = "products-reserved"
)

// Commands.
const (
	// Products.
	CommandTypeAddProduct            = "add-product"
	CommandTypeUpdateProductQuantity = "update-product-quantity"

	// Orders.
	CommandTypeAddOrder        = "add-order"
	CommandTypeReserveProducts = "reserve-products"
)

// Domains.
const (
	DomainProduct = "product"
	DomainOrder   = "order"
)

// Queues.
const (
	QueueProducts        = "products"
	QueueOrders          = "orders"
	QueueReserveProducts = "reserve-products"
)

// MQConnector represents an connector interface to message queue.
type MQConnector interface {
	Close()
	Channel() *amqp.Channel
	Dispather() EventDispatcher
	Start(queueNames []string, consumers []Consumer) error
}

// Consumer provides an interface for queue consumers.
type Consumer interface {
	Handler() func(event Event)
	Queue() string
}

// OrdersStore provides an interface for orders db table.
type OrdersStore interface {
	AddOrder(order Order) error
	GetOrder(id uuid.UUID) (Order, error)
}

// ProductsStore provides an interface for products db table.
type ProductsStore interface {
	AddProduct(product Product) error
	UpdateProductQuantity(product Product) error
	GetProducts() ([]Product, error)
	GetProduct(id uuid.UUID) (Product, error)
	ReserveProducts(orderProducts []OrderProduct) ([]Product, error)
}

// EventStore represents an interface for the event store.
type EventStore interface {
	Apply([]Event) error
	Version() int
	Replay(*int) []Event
}

// Aggregate represents an interface for aggregates.
type Aggregate interface {
	WithID(uuid.UUID) Aggregate
	Handle(Command) (Event, error)
}

// Command represents an interface for commands.
type Command interface {
	Domain() string
	Type() string
	Data() interface{}
}

// Event represents an event model.
type Event struct {
	ID            uuid.UUID       `json:"id"`
	Type          string          `json:"type"`
	AggregateType string          `json:"aggregateType"`
	AggregateID   uuid.UUID       `json:"aggregateId"`
	CreatedAt     time.Time       `json:"created_at"`
	Version       int32           `json:"version"`
	Data          json.RawMessage `json:"data"`
}

// EventDispatcher represents a dispatcher for events.
type EventDispatcher interface {
	Dispatch(queueName string, event Event) error
}

// Product represents a product model.
type Product struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Quantity    int32      `json:"quantity"`
	Price       float64    `json:"price"`
}

// ProductPatch represents a product model with optional properties.
type ProductPatch struct {
	ID       uuid.UUID `json:"id"`
	Quantity *int32    `json:"quantity,omitempty"`
}

//Order represents an order model.
type Order struct {
	ID       *uuid.UUID     `json:"id,omitempty"`
	Products []OrderProduct `json:"products"`
}

// OrderProduct holds data of the order.
type OrderProduct struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Quantity int32     `json:"quantity"`
	Price    *float64  `json:"price,omitempty"`
}
