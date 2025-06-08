package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/s3ndd/sen-go/client"
	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/internal/config"
)

var (
	httpClient     client.ClientInterface
	httpClientSync sync.Once
)

type ServicePrefix string

const (
	RegistryServicePrefix          ServicePrefix = "registry"
	AlertNotificationServicePrefix ServicePrefix = "alert-notification"
	PaymentServicePrefix           ServicePrefix = "payment"
	CartServicePrefix              ServicePrefix = "css"
	UserServicePrefix              ServicePrefix = "user"
	RetailServicePrefix            ServicePrefix = "retail"
	UnifiedSessionServicePrefix    ServicePrefix = "unified_session"
)

// HttpClient returns a singleton rest client
func HttpClient() client.ClientInterface {
	httpClientSync.Do(func() {
		httpClient = client.NewClient(config.HTTPClient())
	})

	return httpClient
}

// Uri returns the full uri.
func Uri(servicePrefix ServicePrefix, version, path string) string {
	return fmt.Sprintf("/%s/%s/%s", servicePrefix, version, path)
}

// GenerateHeaders attaches the api-key to the http header
func GenerateHeaders() map[string]string {
	return map[string]string{
		"api-key": config.ApiKey("SERVICE_API_KEY"),
	}
}

// GenerateSignatureHeaders attaches retailer_id and site_id to the http header
func GenerateSignatureHeaders(retailerID, siteID string) map[string]string {
	return map[string]string{
		"retailer_id": retailerID,
		"site_id":     siteID,
	}
}

// CheckStatus checks the http status and returns the related error message
func CheckStatus(statusCode int) error {
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return nil
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return errors.New("invalid request credentials error")
	case http.StatusForbidden:
		return errors.New("forbidden error")
	case http.StatusNotFound:
		return errors.New("not found error")
	case http.StatusBadRequest:
		return errors.New("bad request error")
	default:
		return errors.New("internal server error")
	}
}

// batchQuery queries the data via the bulk endpoint.
func batchQuery(ctx context.Context, path string, ids []string, response interface{}) error {
	logger := log.ForRequest(ctx).WithFields(log.LogFields{
		"path": path,
		"ids":  ids,
	})
	body, _ := json.Marshal(struct {
		IDs []string `json:"ids"`
	}{
		ids,
	})

	resp, err := HttpClient().Post(ctx, path, GenerateHeaders(), ioutil.NopCloser(bytes.NewBuffer(body)), response)
	if err != nil {
		logger.WithError(err).Error("failed to send the batch query request")
		return err
	}
	if err := CheckStatus(resp.StatusCode()); err != nil {
		logger.WithField("status_code", resp.StatusCode()).WithError(err).
			Error("failed to query the data in batch")
		return err
	}
	logger.WithField("response", response).Info("batch query response")

	return nil
}
