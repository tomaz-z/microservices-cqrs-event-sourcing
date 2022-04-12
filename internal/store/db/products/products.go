package products

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
) escqrs.ProductsStore {
	return store{
		ctx:      ctx,
		dbClient: dbClient,
	}
}

func (s store) AddProduct(product escqrs.Product) error {
	_, err := s.dbClient.PutItem(s.ctx, &dynamodb.PutItemInput{
		TableName: swag.String("Products"),
		Item: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: product.ID.String(),
			},
			"Name": &types.AttributeValueMemberS{
				Value: product.Name,
			},
			"Description": &types.AttributeValueMemberS{
				Value: product.Description,
			},
			"Quantity": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", product.Quantity),
			},
			"Price": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%f", product.Price),
			},
		},
	})
	return err
}

func (s store) UpdateProductQuantity(product escqrs.Product) error {
	_, err := s.GetProduct(*product.ID)
	if err != nil {
		return err
	}

	_, err = s.dbClient.UpdateItem(s.ctx, &dynamodb.UpdateItemInput{
		TableName: swag.String("Products"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: product.ID.String(),
			},
		},
		UpdateExpression: aws.String("set Quantity = :q"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":q": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", product.Quantity),
			},
		},
	})
	return err
}

func (s store) GetProducts() ([]escqrs.Product, error) {
	data := []escqrs.Product{}

	result, err := s.dbClient.Scan(s.ctx, &dynamodb.ScanInput{
		TableName: aws.String("Products"),
	})
	if err != nil {
		return data, err
	}

	for _, item := range result.Items {
		price, err := strconv.ParseFloat(item["Price"].(*types.AttributeValueMemberN).Value, 64)
		if err != nil {
			return data, err
		}
		quantity, err := strconv.Atoi(item["Quantity"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return data, err
		}

		id := uuid.MustParse(item["ID"].(*types.AttributeValueMemberS).Value)

		data = append(data, escqrs.Product{
			ID:          &id,
			Name:        item["Name"].(*types.AttributeValueMemberS).Value,
			Description: item["Description"].(*types.AttributeValueMemberS).Value,
			Price:       price,
			Quantity:    int32(quantity),
		})
	}

	return data, err
}

func (s store) GetProduct(id uuid.UUID) (escqrs.Product, error) {
	data := escqrs.Product{}

	result, err := s.dbClient.GetItem(s.ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Products"),
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

	price, err := strconv.ParseFloat(result.Item["Price"].(*types.AttributeValueMemberN).Value, 64)
	if err != nil {
		return data, err
	}
	quantity, err := strconv.Atoi(result.Item["Quantity"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		return data, err
	}

	resultID := uuid.MustParse(result.Item["ID"].(*types.AttributeValueMemberS).Value)

	return escqrs.Product{
		ID:          &resultID,
		Name:        result.Item["Name"].(*types.AttributeValueMemberS).Value,
		Description: result.Item["Description"].(*types.AttributeValueMemberS).Value,
		Price:       price,
		Quantity:    int32(quantity),
	}, err
}

func (s store) ReserveProducts(orderProducts []escqrs.OrderProduct) ([]escqrs.Product, error) {
	products := []escqrs.Product{}
	for _, item := range orderProducts {
		product, err := s.GetProduct(item.ID)
		if err != nil {
			return []escqrs.Product{}, err
		}

		newQty := product.Quantity - item.Quantity
		if newQty < 0 {
			return []escqrs.Product{}, fmt.Errorf("couldn't reserve %d products of id %s, only %d are left",
				item.Quantity, item.ID.String(), product.Quantity)
		}

		_, err = s.dbClient.UpdateItem(s.ctx, &dynamodb.UpdateItemInput{
			TableName: swag.String("Products"),
			Key: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{
					Value: item.ID.String(),
				},
			},
			UpdateExpression: aws.String("set Quantity = :q"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":q": &types.AttributeValueMemberN{
					Value: fmt.Sprintf("%d", newQty),
				},
			},
		})
		if err != nil {
			return []escqrs.Product{}, err
		}

		products = append(products, product)
	}
	return products, nil
}
