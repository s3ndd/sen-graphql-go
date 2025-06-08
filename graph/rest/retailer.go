package rest

import (
	"context"
	"fmt"

	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/graph/model"
)

func GetRetailerByID(ctx context.Context, retailerID string) (*model.Retailer, error) {
	var retailer model.Retailer
	resp, err := HttpClient().Get(ctx,
		Uri(RegistryServicePrefix, "v1", fmt.Sprintf("retailers/%s?include_sites=true", retailerID)),
		GenerateHeaders(), &retailer)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("retailer_id", retailerID).
			Error("failed to get the retailer with the given id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{"response": retailer, "status_code": resp.StatusCode()}).WithError(err).
			Error("failed to get the retailer from registry service")
		return nil, err
	}
	return &retailer, nil
}

func GetRetailerByIDs(ctx context.Context, retailerIDs []string) ([]*model.Retailer, error) {
	log.Global().WithField("retailer_ids", retailerIDs).Info("get retailers")
	path := Uri(RegistryServicePrefix, "v1", "retailers/bulk")
	var response map[string]*model.Retailer
	if err := batchQuery(ctx, path, retailerIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the retailers by ids from registry service")
		return nil, err
	}
	retailers := make([]*model.Retailer, len(response))
	for i := range retailerIDs {
		if retailer, ok := response[retailerIDs[i]]; ok {
			retailers[i] = retailer
		}
	}
	return retailers, nil
}

func GetRetailers(ctx context.Context) ([]*model.Retailer, error) {
	var retailers struct {
		Retailers []*model.Retailer `json:"retailers"`
	}
	resp, err := HttpClient().Get(ctx,
		Uri(RegistryServicePrefix, "v1", "retailers"),
		GenerateHeaders(), &retailers)
	if err != nil {
		log.ForRequest(ctx).WithError(err).Error("failed to get the retailer with the given id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{"response": retailers, "status_code": resp.StatusCode()}).WithError(err).
			Error("failed to get the retailer from registry service")
		return nil, err
	}
	return retailers.Retailers, nil
}
