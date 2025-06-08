package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newSessionLoader(ctx context.Context) *SessionLoader {
	return &SessionLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(sessionIDs []string) ([]*model.Session, []error) {
			sessions, err := rest.GetSessionByIDs(ctx, sessionIDs)
			return sessions, errors(err)
		},
	}
}
