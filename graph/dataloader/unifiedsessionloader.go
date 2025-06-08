package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newUnifiedSessionLoader(ctx context.Context) *UnifiedSessionLoader {
	return &UnifiedSessionLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(sessionIDs []string) ([]*model.Session, []error) {
			sessions, err := rest.GetUnifiedSessionByIDs(ctx, sessionIDs)
			return sessions, errors(err)
		},
	}
}
