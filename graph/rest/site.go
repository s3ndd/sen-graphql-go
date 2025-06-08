package rest

import (
	"context"
	"fmt"

	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/graph/model"
)

func GetSitesByRetailerID(ctx context.Context, retailerID string) ([]*model.Site, error) {
	var sites struct {
		Sites []*model.Site `json:"sites"`
	}
	resp, err := HttpClient().Get(ctx,
		Uri(RegistryServicePrefix, "v1", fmt.Sprintf("retailers/%s/sites", retailerID)),
		GenerateHeaders(), &sites)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("retailer_id", retailerID).
			Error("failed to get the sites with the given retailer id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{"response": sites, "status_code": resp.StatusCode()}).WithError(err).
			Error("failed to get the sites from registry service")
		return nil, err
	}
	return sites.Sites, nil
}

func GetSiteByID(ctx context.Context, siteID, retailerID string) (*model.Site, error) {
	var site model.Site
	resp, err := HttpClient().Get(ctx,
		Uri(RegistryServicePrefix, "v1", fmt.Sprintf("retailers/%s/sites/%s?hide_sensitive=false", retailerID, siteID)),
		GenerateHeaders(), &site)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("retailer_id", retailerID).
			Error("failed to get the sites with the given retailer id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{"response": site, "status_code": resp.StatusCode()}).WithError(err).
			Error("failed to get the sites from registry service")
		return nil, err
	}
	return &site, nil
}

func GetSiteByIDs(ctx context.Context, siteIDs []string) ([]*model.Site, error) {
	path := Uri(RegistryServicePrefix, "v1", "sites/bulk")
	var response map[string]*model.Site
	if err := batchQuery(ctx, path, siteIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the sites by ids from registry service")
		return nil, err
	}
	sites := make([]*model.Site, len(siteIDs))
	for i := range siteIDs {
		if site, ok := response[siteIDs[i]]; ok {
			sites[i] = site
		}
	}

	return sites, nil
}
