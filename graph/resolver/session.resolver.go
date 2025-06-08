package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"github.com/s3ndd/sen-graphql-go/graph/dataloader"
	"github.com/s3ndd/sen-graphql-go/graph/generated"
	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/resolver/helper"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
	"strings"

	"github.com/s3ndd/sen-go/log"
)

func (r *mutationResolver) UpdateSessionStatus(ctx context.Context, input *model.UpdateSessionStatusRequest) (*model.Session, error) {
	logger := log.ForRequest(ctx).WithFields(log.LogFields{
		"session_id": input.SessionID,
		"status":     input.Status,
	})
	session, err := helper.LoadSessionByID(ctx, input.SessionID)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the session loader")
		return nil, err
	}
	if session == nil {
		err = fmt.Errorf("failed to get session by id")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	session, err = rest.UpdateSessionById(ctx, input)
	if err != nil {
		logger.WithError(err).Error("failed to update session by id")
		return nil, err
	}
	return session, nil
}

func (r *queryResolver) Session(ctx context.Context, id string, siteID string, retailerID string) (*model.Session, error) {
	logger := log.ForRequest(ctx).WithFields(log.LogFields{
		"session_id":  id,
		"site_id":     siteID,
		"retailer_id": retailerID,
	})
	loaders := dataloader.ContextLoaders(ctx)
	session, err := loaders.Sessions.Load(id)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the session loader")
		return nil, err
	}
	if session != nil {
		if session.SiteID != siteID || (session.RetailerID != nil && *session.RetailerID != retailerID) {
			logger.Warn("session with the given id is not found")
			return nil, errors.New("session with the given id is not found")
		}
		// attach the retailer id to session
		session.RetailerID = &retailerID
		return session, nil
	}

	session, err = loaders.UnifiedSessions.Load(id)
	if err != nil {
		logger.WithError(err).Error("failed to get the session from the unified session loader")
		return nil, err
	}
	if session == nil || session.SiteID != siteID || (session.RetailerID != nil && *session.RetailerID != retailerID) {
		return nil, errors.New("session with the given id is not found")
	}
	return session, nil
}

func (r *queryResolver) Sessions(ctx context.Context, siteID string, retailerID string, status []model.SessionStatus) (*model.SessionConnection, error) {
	sessions, err := rest.GetSessionsBySiteID(ctx, siteID, status)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionResolver) RepeatUser(ctx context.Context, obj *model.Session) (*bool, error) {
	repeatUser := false
	if obj.UserID == nil {
		return &repeatUser, nil
	}
	sessionConnection, err := rest.GetSessionsByUserID(ctx, *obj.UserID, 1, 2)
	if err != nil {
		log.ForRequest(ctx).WithField("session", obj).WithError(err).
			Warn("failed to retrieve sessions by user id, so return false by default")
		return &repeatUser, nil
	}
	repeatUser = len(sessionConnection.Sessions) > 1
	return &repeatUser, nil
}

func (r *sessionResolver) Site(ctx context.Context, obj *model.Session) (*model.Site, error) {
	loaders := dataloader.ContextLoaders(ctx)
	site, err := loaders.Sites.Load(obj.SiteID)
	if err != nil || site == nil {
		return nil, err
	}
	return site, nil
}

func (r *sessionResolver) Status(ctx context.Context, obj *model.Session) (model.SessionStatus, error) {
	if strings.HasPrefix(obj.Status, "SessionStatus_") {
		return model.SessionStatus(strings.ReplaceAll(obj.Status, "SessionStatus_", "")), nil
	}
	return model.SessionStatus(obj.Status), nil
}

func (r *sessionResolver) Alerts(ctx context.Context, obj *model.Session, status []model.AlertStatus, types []model.AlertType) ([]*model.Alert, error) {
	alerts, err := rest.GetAlertsBySessionID(ctx, obj.ID, status, types)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *sessionResolver) Events(ctx context.Context, obj *model.Session, eventType *model.EventType, eventSubTypes []model.EventSubType) ([]*model.Event, error) {
	if obj.RetailerID == nil {
		loaders := dataloader.ContextLoaders(ctx)
		site, err := loaders.Sites.Load(obj.SiteID)
		if err != nil || site == nil {
			return nil, err
		}
		obj.RetailerID = &site.RetailerID
	}
	eventsResponse, err := rest.GetEventsBySessionID(ctx, obj.ID, obj.SiteID, *obj.RetailerID, eventType, eventSubTypes)
	if err != nil {
		return nil, err
	}
	log.ForRequest(ctx).WithFields(log.LogFields{
		"session_id":      obj.ID,
		"retailer_id":     obj.RetailerID,
		"site_id":         obj.SiteID,
		"event_type":      eventType,
		"event_sub_types": eventSubTypes,
	}).WithField("events", eventsResponse).Debug("get events by session id response")

	events := make([]*model.Event, 0, len(eventsResponse))

	for _, eventResponse := range eventsResponse {
		event := &model.Event{
			ID:           eventResponse.ID,
			SessionID:    eventResponse.SessionID,
			EventType:    eventResponse.EventType,
			EventSubType: eventResponse.Message.EventSubType,
			Flagged:      eventResponse.Flagged,
			Skipped:      eventResponse.Skipped,
			Created:      eventResponse.Created,
			Updated:      eventResponse.Updated,
		}
		// compatible with retail service response
		if len(eventResponse.Message.Items) > 0 {
			event.ProductKeyList = eventResponse.Message.Items
		}
		if len(event.ProductKeyList) == 0 {
			event.ProductKeyList = eventResponse.Message.ProductKeyList
			quantityOne := 1
			// this quantity adjust is for the retail service only
			for i := range event.ProductKeyList {
				if event.ProductKeyList[i].Quantity == nil {
					event.ProductKeyList[i].Quantity = &quantityOne
				}
			}
		}

		events = append(events, event)
	}

	return events, nil
}

// Session returns generated.SessionResolver implementation.
func (r *Resolver) Session() generated.SessionResolver { return &sessionResolver{r} }

type sessionResolver struct{ *Resolver }
