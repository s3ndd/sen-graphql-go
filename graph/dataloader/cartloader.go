package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newCartLoader(ctx context.Context) *CartLoader {
	return &CartLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(cartIDs []string) ([]*model.Cart, []error) {
			carts, err := rest.GetCartByIDs(ctx, cartIDs)
			return carts, errors(err)
		},
	}
}
