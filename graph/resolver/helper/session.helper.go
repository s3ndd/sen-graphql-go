package helper

import (
	"context"
	"github.com/s3ndd/sen-go/log"
	"github.com/s3ndd/sen-graphql-go/graph/dataloader"
	"github.com/s3ndd/sen-graphql-go/graph/model"
)

func LoadSessionByID(ctx context.Context, id string) (*model.Session, error) {
	logger := log.ForRequest(ctx).WithField("session_id", id)
	loaders := dataloader.ContextLoaders(ctx)
	session, err := loaders.Sessions.Load(id)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the session loader")
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	session, err = loaders.UnifiedSessions.Load(id)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the unified session loader")
		return nil, err
	}
	return session, nil
}
