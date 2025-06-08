package rest

import (
	"context"
	"fmt"

	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/graph/model"
)

// GetCartsRegistrationBySiteID returns the cart from registry service
func GetCartsRegistrationBySiteID(ctx context.Context, siteID string) (*model.CartsRegistration, error) {
	var carts model.CartsRegistration
	resp, err := HttpClient().Get(ctx,
		// TODO: include_deleted=false&page_size=100 will be updated later
		Uri(RegistryServicePrefix, "v1", fmt.Sprintf("carts?site_id=%s&include_deleted=false&page_size=100", siteID)),
		GenerateHeaders(), &carts)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("site_id", siteID).
			Error("failed to get the carts with the given site id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    carts,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the carts with the given site id from registry service")
		return nil, err
	}
	return &carts, nil
}

// GetCartRegistrationByID returns the cart from registry service
func GetCartRegistrationByID(ctx context.Context, id string) (*model.CartRegistration, error) {
	var cart model.CartRegistration
	resp, err := HttpClient().Get(ctx,
		Uri(RegistryServicePrefix, "v1", fmt.Sprintf("carts/id/%s", id)),
		GenerateHeaders(), &cart)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("cart_id", id).
			Error("failed to get the cart with the given id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    cart,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the cart with the given id from registry service")
		return nil, err
	}
	return &cart, nil
}

func GetCartRegistrationByIDs(ctx context.Context, cartIDs []string) ([]*model.CartRegistration, error) {
	path := Uri(RegistryServicePrefix, "v1", "carts/bulk")
	var response map[string]*model.CartRegistration
	if err := batchQuery(ctx, path, cartIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the carts by ids from registry service")
		return nil, err
	}
	carts := make([]*model.CartRegistration, len(response))
	for i := range cartIDs {
		if cart, ok := response[cartIDs[i]]; ok {
			carts[i] = cart
		}
	}

	return carts, nil
}

// GetCartRegistrationByQRCode returns the cart from registry service
func GetCartRegistrationByQRCode(ctx context.Context, qrCode string) (*model.CartRegistration, error) {
	var cart model.CartRegistration
	resp, err := HttpClient().Get(ctx,
		Uri(RegistryServicePrefix, "v1", fmt.Sprintf("carts/qr_code/%s", qrCode)),
		GenerateHeaders(), &cart)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("cart_qr_code", qrCode).
			Error("failed to get the cart with the given qr code")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    cart,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the cart with the given qr code from registry service")
		return nil, err
	}
	return &cart, nil
}

func GetCartByID(ctx context.Context, id, siteID string) (*model.Cart, error) {
	var cart model.Cart
	resp, err := HttpClient().Get(ctx,
		Uri(CartServicePrefix,
			"v1",
			fmt.Sprintf("%s/carts/id/%s?embed=%s&embed=%s&embed=%s", siteID, id, "wifi", "temperature", "battery")),
		GenerateHeaders(), &cart)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithFields(log.LogFields{
			"cart_id": id,
			"site_id": siteID,
		}).Error("failed to get the cart with the given id from css")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    cart,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the cart with the given id from css")
		return nil, err
	}
	return &cart, nil
}

func GetCartByIDs(ctx context.Context, cartIDs []string) ([]*model.Cart, error) {
	path := Uri(CartServicePrefix, "v1", "carts/bulk")
	var response map[string]*model.Cart
	if err := batchQuery(ctx, path, cartIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the carts by ids from cart service")
		return nil, err
	}
	carts := make([]*model.Cart, len(response))
	for i := range cartIDs {
		if cart, ok := response[cartIDs[i]]; ok {
			carts[i] = cart
		}
	}

	return carts, nil
}

func GetCartByQRCode(ctx context.Context, qrCode, siteID string) (*model.Cart, error) {
	var cart model.Cart
	resp, err := HttpClient().Get(ctx,
		Uri(CartServicePrefix,
			"v1",
			fmt.Sprintf("%s/carts/code/%s?embed=%s&embed=%s&embed=%s", siteID, qrCode, "wifi", "temperature", "battery")),
		GenerateHeaders(), &cart)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithFields(log.LogFields{
			"cart_qr_code": qrCode,
			"site_id":      siteID,
		}).Error("failed to get the cart with the given qr code from css")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    cart,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the cart with the given qr code from css")
		return nil, err
	}
	return &cart, nil
}

func GetCartsBySiteID(ctx context.Context, siteID string) (*model.CartConnection, error) {
	var carts model.CartConnection
	resp, err := HttpClient().Get(ctx,
		Uri(CartServicePrefix,
			"v1",
			fmt.Sprintf("%s/carts?embed=%s&embed=%s&embed=%s", siteID, "wifi", "temperature", "battery")),
		GenerateHeaders(), &carts)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("site_id", siteID).
			Error("failed to get the carts with the given site id from css")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    carts,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the carts with the given site id from css")
		return nil, err
	}
	return &carts, nil
}
