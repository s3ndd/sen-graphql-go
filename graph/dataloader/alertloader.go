package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newAlertLoader(ctx context.Context) *AlertLoader {
	return &AlertLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(cartIDs []string) ([]*model.Alert, []error) {
			alerts, err := rest.GetAlertByIDs(ctx, cartIDs)
			return alerts, errors(err)
		},
	}
}
