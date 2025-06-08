package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/s3ndd/sen-graphql-go/graph/dataloader"
	"github.com/s3ndd/sen-graphql-go/graph/generated"
	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func (r *alertResolver) Session(ctx context.Context, obj *model.Alert) (*model.Session, error) {
	if obj.SessionID == nil {
		return nil, nil
	}
	loaders := dataloader.ContextLoaders(ctx)
	return loaders.Sessions.Load(*obj.SessionID)
}

func (r *queryResolver) Alert(ctx context.Context, id string) (*model.Alert, error) {
	loaders := dataloader.ContextLoaders(ctx)
	return loaders.Alerts.Load(id)
}

func (r *queryResolver) Alerts(ctx context.Context, siteID *string, sessionID *string, cartID *string, status []model.AlertStatus, types []model.AlertType) (*model.AlertConnection, error) {
	alerts, err := rest.GetAlerts(ctx, siteID, sessionID, cartID, status, types)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// Alert returns generated.AlertResolver implementation.
func (r *Resolver) Alert() generated.AlertResolver { return &alertResolver{r} }

type alertResolver struct{ *Resolver }
