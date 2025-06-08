package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newCartRegistrationLoader(ctx context.Context) *CartRegistrationLoader {
	return &CartRegistrationLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(cartIDs []string) ([]*model.CartRegistration, []error) {
			carts, err := rest.GetCartRegistrationByIDs(ctx, cartIDs)
			return carts, errors(err)
		},
	}
}
