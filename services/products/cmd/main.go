package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/go-openapi/loads"

	conProducts "eventsourcing/internal/consumers/products"
	conReserveProducts "eventsourcing/internal/consumers/reserve_products"
	"eventsourcing/internal/mq"
	"eventsourcing/internal/store/db"
	productsStore "eventsourcing/internal/store/db/products"
	"eventsourcing/internal/store/events"
	escqrs "eventsourcing/services"
	productsAddH "eventsourcing/services/products/api/handlers/products/add"
	productsPatchH "eventsourcing/services/products/api/handlers/products/patch"
	productsProductH "eventsourcing/services/products/api/handlers/products/product"
	productsProductsH "eventsourcing/services/products/api/handlers/products/products"
	"eventsourcing/services/products/api/restapi"
	"eventsourcing/services/products/api/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewProductsAPI(swaggerSpec)
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
	dbProducts, err := db.New(ctx, host)
	if err != nil {
		log.Fatal(err)
	}
	dbEvents, err := db.New(ctx, host_events)
	if err != nil {
		log.Fatal(err)
	}

	productsStorage := productsStore.New(ctx, dbProducts)
	eventStore := events.New(
		ctx,
		dbEvents,
		"Products",
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
			escqrs.QueueProducts,
			escqrs.QueueReserveProducts,
		},
		[]escqrs.Consumer{
			conProducts.New(
				eventStore,
				productsStorage,
			),
			conReserveProducts.New(
				eventStore,
				productsStorage,
				mqDispatcher,
			),
		},
	)
	if err != nil {
		log.Fatal("starting queues and consumers: ", err)
	}

	// Command handlers.
	api.ProductsAddHandler = productsAddH.New(
		eventStore,
		mqDispatcher,
	)
	api.ProductsPatchHandler = productsPatchH.New(
		eventStore,
		mqDispatcher,
	)

	// Query handlers.
	api.ProductsProductsHandler = productsProductsH.New(
		productsStorage,
	)
	api.ProductsProductHandler = productsProductH.New(
		productsStorage,
	)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
