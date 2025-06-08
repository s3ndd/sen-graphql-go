package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/resolver/helper"
	"github.com/s3ndd/sen-graphql-go/graph/rest"

	"github.com/s3ndd/sen-go/log"
)

func (r *mutationResolver) AddItems(ctx context.Context, input *model.ItemsRequest) (*model.EventResponse, error) {
	logger := log.ForRequest(ctx).WithField("input", input)
	// verify session id, and obtain site id and retailer id
	session, err := helper.LoadSessionByID(ctx, input.SessionID)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the session loader")
		return nil, err
	}
	if session == nil {
		err = fmt.Errorf("failed to get the session by id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	if session.SiteID != input.SiteID {
		err = fmt.Errorf("session not found with the given site id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	event, err := rest.ProcessItems(ctx, input, model.AddActionType)
	if err != nil {
		logger.WithError(err).Error("failed to process add item event")
		return nil, err
	}
	return event, nil
}

func (r *mutationResolver) RemoveItems(ctx context.Context, input *model.ItemsRequest) (*model.EventResponse, error) {
	logger := log.ForRequest(ctx).WithField("input", input)
	// verify session id, and obtain site id and retailer id
	session, err := helper.LoadSessionByID(ctx, input.SessionID)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the session loader")
		return nil, err
	}
	if session == nil {
		err = fmt.Errorf("failed to get the session by id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	if session.SiteID != input.SiteID {
		err = fmt.Errorf("session not found with the given site id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	event, err := rest.ProcessItems(ctx, input, model.RemoveActionType)
	if err != nil {
		logger.WithError(err).Error("failed to process remove item event")
		return nil, err
	}
	return event, nil
}

func (r *mutationResolver) ReplaceItem(ctx context.Context, input *model.ReplaceItemRequest) (*model.EventResponse, error) {
	logger := log.ForRequest(ctx).WithField("input", input)
	session, err := helper.LoadSessionByID(ctx, input.SessionID)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the session loader")
		return nil, err
	}
	if session == nil {
		err = fmt.Errorf("failed to get the session by id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	if session.SiteID != input.SiteID {
		err = fmt.Errorf("session not found with the given site id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	fromKey, toKey := input.FromItem.ProductKey, input.ToItem.ProductKey
	if input.FromItem.Discount != nil {
		fromKey += *input.FromItem.Discount
	}
	if input.ToItem.Discount != nil {
		toKey += *input.ToItem.Discount
	}
	if fromKey == toKey {
		err = fmt.Errorf("replacement item cannot be identical")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	event, err := rest.ReplaceItem(ctx, input)
	return event, err
}
