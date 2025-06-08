package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/s3ndd/sen-graphql-go/graph/dataloader"
	"github.com/s3ndd/sen-graphql-go/graph/generated"
	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/resolver/helper"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func (r *cartResolver) Site(ctx context.Context, obj *model.Cart) (*model.Site, error) {
	loaders := dataloader.ContextLoaders(ctx)
	cart, err := loaders.CartRegistrations.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	return cart.Site, nil
}

func (r *cartResolver) Session(ctx context.Context, obj *model.Cart) (*model.Session, error) {
	if obj.SessionID == nil {
		return nil, nil
	}

	return helper.LoadSessionByID(ctx, *obj.SessionID)
}

func (r *cartResolver) Sessions(ctx context.Context, obj *model.Cart, status []model.SessionStatus) ([]*model.Session, error) {
	sessions, err := rest.GetSessionsBySiteID(ctx, obj.SiteID, status)
	if err != nil {
		return nil, err
	}
	return sessions.Sessions, nil
}

func (r *cartResolver) Alerts(ctx context.Context, obj *model.Cart, status []model.AlertStatus, types []model.AlertType) ([]*model.Alert, error) {
	alerts, err := rest.GetAlerts(ctx, &obj.SiteID, obj.SessionID, &obj.ID, status, types)
	if err != nil {
		return nil, err
	}

	return alerts.Alerts, nil
}

func (r *queryResolver) Cart(ctx context.Context, id string) (*model.Cart, error) {
	loaders := dataloader.ContextLoaders(ctx)
	return loaders.Carts.Load(id)
}

func (r *queryResolver) Carts(ctx context.Context, siteID string) (*model.CartConnection, error) {
	cartConnection, err := rest.GetCartsBySiteID(ctx, siteID)
	if err != nil {
		return nil, err
	}

	return cartConnection, nil
}

// Cart returns generated.CartResolver implementation.
func (r *Resolver) Cart() generated.CartResolver { return &cartResolver{r} }

type cartResolver struct{ *Resolver }
