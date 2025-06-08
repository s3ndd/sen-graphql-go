package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newRetailerLoader(ctx context.Context) *RetailerLoader {
	return &RetailerLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(retailerIDs []string) ([]*model.Retailer, []error) {
			retailers, err := rest.GetRetailerByIDs(ctx, retailerIDs)
			return retailers, errors(err)
		},
	}
}
