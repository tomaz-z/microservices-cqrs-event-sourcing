package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/go-openapi/loads"

	conOrders "eventsourcing/internal/consumers/orders"
	"eventsourcing/internal/mq"
	"eventsourcing/internal/store/db"
	"eventsourcing/internal/store/db/orders"
	"eventsourcing/internal/store/events"
	escqrs "eventsourcing/services"
	ordersCreateH "eventsourcing/services/orders/api/handlers/orders/create"
	ordersOrderH "eventsourcing/services/orders/api/handlers/orders/order"
	"eventsourcing/services/orders/api/restapi"
	"eventsourcing/services/orders/api/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewOrdersAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	flag.Parse()
	server.Port = 80

	ctx := context.Background()

	host := os.Getenv("DB_HOST")
	if len(host) < 1 {
		log.Fatal("couldn't read DB_HOST")
	}

	host_events := os.Getenv("DB_EVENTS_HOST")
	if len(host_events) < 1 {
		log.Fatal("couldn't read DB_EVENTS_HOST")
	}

	// Prerequisites.
	dbOrders, err := db.New(ctx, host)
	if err != nil {
		log.Fatal(err)
	}
	dbEvents, err := db.New(ctx, host_events)
	if err != nil {
		log.Fatal(err)
	}

	ordersStore := orders.New(ctx, dbOrders)
	eventStore := events.New(
		ctx,
		dbEvents,
		"Orders",
	)

	// Message queue.
	connector, err := mq.Connect("amqp://admin:admin@mq:5672/")
	if err != nil {
		log.Fatal("connecting to MQ: ", err)
	}
	defer connector.Close()
	mqDispatcher := connector.Dispather()

	err = connector.Start(
		[]string{
			escqrs.QueueOrders,
			escqrs.QueueReserveProducts,
		},
		[]escqrs.Consumer{
			conOrders.New(
				eventStore,
				ordersStore,
			),
		},
	)
	if err != nil {
		log.Fatal("starting queues and consumers: ", err)
	}

	// Command handlers.
	api.OrdersCreateHandler = ordersCreateH.New(mqDispatcher)

	// Query handlers.
	api.OrdersOrderHandler = ordersOrderH.New(ordersStore)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
