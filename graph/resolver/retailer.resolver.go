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

func (r *queryResolver) Retailer(ctx context.Context, id string) (*model.Retailer, error) {
	loaders := dataloader.ContextLoaders(ctx)
	return loaders.Retailers.Load(id)
}

func (r *queryResolver) Retailers(ctx context.Context) ([]*model.Retailer, error) {
	retailers, err := rest.GetRetailers(ctx)
	if err != nil {
		return nil, err
	}
	return retailers, nil
}

func (r *retailerResolver) Sites(ctx context.Context, obj *model.Retailer) ([]*model.Site, error) {
	sites, err := rest.GetSitesByRetailerID(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

// Retailer returns generated.RetailerResolver implementation.
func (r *Resolver) Retailer() generated.RetailerResolver { return &retailerResolver{r} }

type retailerResolver struct{ *Resolver }
