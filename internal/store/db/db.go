package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// New returns a new dynamodb client instance.
func New(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	config, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolver(resolver{
			endpoint: endpoint,
		}),
	)
	if err != nil {
		return nil, err
	}
	config.Retryer = func() aws.Retryer {
		standard := retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = 20
			o.RateLimiter = noopRateLimit{}
		})
		return standard
	}

	return dynamodb.NewFromConfig(config), nil
}

type noopRateLimit struct{}

func (n noopRateLimit) GetToken(ctx context.Context, cost uint) (releaseToken func() error, err error) {
	return func() error { return nil }, nil
}
func (n noopRateLimit) AddTokens(uint) error { return nil }

type resolver struct {
	endpoint string
}

func (r resolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	if region == "local" {
		return aws.Endpoint{
			URL: r.endpoint,
		}, nil
	}
	return aws.Endpoint{}, errors.New("unknown endpoint requested")
}
