package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/s3ndd/sen-graphql-go/graph/generated"
	"github.com/s3ndd/sen-graphql-go/graph/model"
)

func (r *itemResolver) DiscountType(ctx context.Context, obj *model.Item) (string, error) {
	return obj.DiscountType.String(), nil
}

func (r *queryResolver) Items(ctx context.Context, sessionID string) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

// Item returns generated.ItemResolver implementation.
func (r *Resolver) Item() generated.ItemResolver { return &itemResolver{r} }

type itemResolver struct{ *Resolver }
