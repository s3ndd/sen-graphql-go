package dataloader

import (
	"context"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/graph/rest"
)

func newSiteLoader(ctx context.Context) *SiteLoader {
	return &SiteLoader{
		wait:     waitTime,
		maxBatch: maxBatch,
		fetch: func(siteIDs []string) ([]*model.Site, []error) {
			sites, err := rest.GetSiteByIDs(ctx, siteIDs)
			return sites, errors(err)
		},
	}
}
