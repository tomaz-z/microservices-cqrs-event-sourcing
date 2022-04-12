package events

import (
	"context"
	escqrs "eventsourcing/services"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type eventStore struct {
	ctx      context.Context
	dbClient *dynamodb.Client
	domain   string
	Events   []escqrs.Event
}

// New returns a new event store.
func New(
	ctx context.Context,
	dbClient *dynamodb.Client,
	domain string,
) escqrs.EventStore {
	return &eventStore{
		ctx:      ctx,
		dbClient: dbClient,
		domain:   "Events",
	}
}

func (e *eventStore) Version() int {
	return len(e.Events)
}

func (e *eventStore) Apply(events []escqrs.Event) error {
	writeRequests := []types.WriteRequest{}
	for _, event := range events {
		result, err := e.dbClient.Scan(e.ctx, &dynamodb.ScanInput{
			TableName: aws.String(e.domain),
		})
		if err != nil {
			return err
		}

		event.Version = result.Count + 1

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"ID": &types.AttributeValueMemberS{
						Value: event.ID.String(),
					},
					"Type": &types.AttributeValueMemberS{
						Value: event.Type,
					},
					"CreatedAt": &types.AttributeValueMemberS{
						Value: event.CreatedAt.String(),
					},
					"Version": &types.AttributeValueMemberN{
						Value: fmt.Sprintf("%d", event.Version),
					},
					"Data": &types.AttributeValueMemberS{
						Value: string(event.Data),
					},
				},
			},
		})
	}

	_, err := e.dbClient.BatchWriteItem(e.ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			e.domain: writeRequests,
		},
	})
	return err
}

func (e *eventStore) Replay(version *int) []escqrs.Event {
	if version == nil {
		return e.Events
	}

	for i, item := range e.Events {
		if item.Version == int32(*version) {
			return e.Events[i:]
		}
	}

	return []escqrs.Event{}
}
