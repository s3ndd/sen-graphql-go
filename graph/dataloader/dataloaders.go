package dataloader

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/s3ndd/sen-graphql-go/internal/middleware"
)

//go:generate go run github.com/vektah/dataloaden RetailerLoader string *github.com/s3ndd/sen-graphql-go/graph/model.Retailer
//go:generate go run github.com/vektah/dataloaden SiteLoader string *github.com/s3ndd/sen-graphql-go/graph/model.Site
//go:generate go run github.com/vektah/dataloaden CartLoader string *github.com/s3ndd/sen-graphql-go/graph/model.Cart
//go:generate go run github.com/vektah/dataloaden CartRegistrationLoader string *github.com/s3ndd/sen-graphql-go/graph/model.CartRegistration
//go:generate go run github.com/vektah/dataloaden AlertLoader string *github.com/s3ndd/sen-graphql-go/graph/model.Alert
//go:generate go run github.com/vektah/dataloaden SessionLoader string *github.com/s3ndd/sen-graphql-go/graph/model.Session
//go:generate go run github.com/vektah/dataloaden UnifiedSessionLoader string *github.com/s3ndd/sen-graphql-go/graph/model.Session

const (
	maxBatch = 1000
	waitTime = 1 * time.Millisecond
)

const ctxKey = "data_loaders"

type Loaders struct {
	Retailers         *RetailerLoader
	Sites             *SiteLoader
	Carts             *CartLoader
	CartRegistrations *CartRegistrationLoader
	Alerts            *AlertLoader
	Sessions          *SessionLoader
	UnifiedSessions   *UnifiedSessionLoader
}

func LoaderMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loaders Loaders

		loaders.Retailers = newRetailerLoader(ctx)
		loaders.Sites = newSiteLoader(ctx)
		loaders.Carts = newCartLoader(ctx)
		loaders.CartRegistrations = newCartRegistrationLoader(ctx)
		loaders.Alerts = newAlertLoader(ctx)
		loaders.Sessions = newSessionLoader(ctx)
		loaders.UnifiedSessions = newUnifiedSessionLoader(ctx)
		ctx.Set(ctxKey, &loaders)
		ctx.Next()
	}
}

func ContextLoaders(ctx context.Context) *Loaders {
	// both context switch and loader will be used as a middleware, so it must be in the context
	// ignore any error here.
	gc, _ := middleware.GinContextFromContext(ctx)
	value, _ := gc.Get(ctxKey)
	return value.(*Loaders)
}

func errors(err error) []error {
	if err != nil {
		return []error{err}
	}
	return nil
}
