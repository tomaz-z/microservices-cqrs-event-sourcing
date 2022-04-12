package orders

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	escqrs "eventsourcing/services"
)

type store struct {
	ctx      context.Context
	dbClient *dynamodb.Client
}

// New returns a new instance of domainvent processor for product domain.
func New(
	ctx context.Context,
	dbClient *dynamodb.Client,
) escqrs.OrdersStore {
	return store{
		ctx:      ctx,
		dbClient: dbClient,
	}
}

func (s store) AddOrder(order escqrs.Order) error {
	products := []types.AttributeValue{}
	for _, item := range order.Products {
		products = append(products, &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{
					Value: item.ID.String(),
				},
				"Quantity": &types.AttributeValueMemberN{
					Value: fmt.Sprintf("%d", item.Quantity),
				},
				"Price": &types.AttributeValueMemberN{
					Value: fmt.Sprintf("%f", *item.Price),
				},
			},
		})
	}

	_, err := s.dbClient.PutItem(s.ctx, &dynamodb.PutItemInput{
		TableName: swag.String("Orders"),
		Item: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: order.ID.String(),
			},
			"Products": &types.AttributeValueMemberL{
				Value: products,
			},
		},
	})
	return err
}

func (s store) GetOrder(id uuid.UUID) (escqrs.Order, error) {
	data := escqrs.Order{}

	result, err := s.dbClient.GetItem(s.ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Orders"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: id.String(),
			},
		},
	})
	if err != nil {
		return data, err
	}

	if len(result.Item) < 1 {
		return data, errors.New("no item found")
	}

	products := []escqrs.OrderProduct{}
	for _, item := range result.Item["Products"].(*types.AttributeValueMemberL).Value {
		values := item.(*types.AttributeValueMemberM).Value

		price, err := strconv.ParseFloat(values["Price"].(*types.AttributeValueMemberN).Value, 64)
		if err != nil {
			return data, err
		}
		quantity, err := strconv.Atoi(values["Quantity"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return data, err
		}

		productID := uuid.MustParse(values["ID"].(*types.AttributeValueMemberS).Value)

		products = append(products, escqrs.OrderProduct{
			ID:       productID,
			Price:    &price,
			Quantity: int32(quantity),
		})
	}

	return escqrs.Order{
		ID:       &id,
		Products: products,
	}, err
}
