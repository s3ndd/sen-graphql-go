package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/s3ndd/sen-graphql-go/graph/dataloader"
	"github.com/s3ndd/sen-graphql-go/graph/generated"
	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func (r *queryResolver) Site(ctx context.Context, retailerID string, id string) (*model.Site, error) {
	loaders := dataloader.ContextLoaders(ctx)
	site, err := loaders.Sites.Load(id)
	if err != nil {
		return nil, err
	}
	if site != nil && site.RetailerID != retailerID {
		return nil, errors.New("site not found")
	}
	return site, nil
}

func (r *queryResolver) Sites(ctx context.Context, retailerID string) ([]*model.Site, error) {
	sites, err := rest.GetSitesByRetailerID(ctx, retailerID)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *siteResolver) IntegrationType(ctx context.Context, obj *model.Site) (string, error) {
	return obj.IntegrationType.String(), nil
}

func (r *siteResolver) Sessions(ctx context.Context, obj *model.Site, status []model.SessionStatus) (*model.SessionConnection, error) {
	sessions, err := rest.GetSessionsBySiteID(ctx, obj.ID, status)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *siteResolver) Alerts(ctx context.Context, obj *model.Site, status []model.AlertStatus, types []model.AlertType) (*model.AlertConnection, error) {
	alerts, err := rest.GetAlerts(ctx, &obj.ID, nil, nil, status, types)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// Site returns generated.SiteResolver implementation.
func (r *Resolver) Site() generated.SiteResolver { return &siteResolver{r} }

type siteResolver struct{ *Resolver }
