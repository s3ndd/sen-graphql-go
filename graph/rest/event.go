package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/s3ndd/sen-go/auth"
	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/internal/middleware"
)

type ShoppingEventRequest struct {
	path    string
	body    []byte
	headers map[string]string
}

type CreateEventBody struct {
	ID        string             `json:"id"`
	SessionID string             `json:"session_id"`
	EventType model.EventType    `json:"event_type"`
	Message   model.EventMessage `json:"message"`
	Flagged   bool               `json:"flagged"`
}

func ProcessItems(ctx context.Context, input *model.ItemsRequest, actionType model.ActionType) (*model.EventResponse, error) {
	logger := log.ForRequest(ctx).WithFields(log.LogFields{"input": input, "action": actionType})
	site, err := GetSiteByID(ctx, input.SiteID, input.RetailerID)
	if err != nil {
		logger.WithError(err).Error("failed to get site by id")
		return nil, err
	}
	body := &CreateEventBody{
		ID:        input.ID,
		SessionID: input.SessionID,
		EventType: model.EventTypeShopping,
		Flagged:   input.Flagged,
		Message: model.EventMessage{
			EventSubType:   actionType.ToEventSubType(),
			WebhookNotify:  input.WebhookNotify,
			ProductKeyList: input.Items,
		},
	}

	request := buildShoppingEventRequest(ctx, site, body)

	var event *model.EventResponse
	resp, err := HttpClient().Post(ctx, request.path, request.headers, ioutil.NopCloser(bytes.NewBuffer(request.body)), &event)
	if err != nil {
		logger.WithError(err).Error(fmt.Sprintf("failed to post %s item event request", actionType))
		return nil, err
	}
	logger.WithField("resp_status_code", resp.StatusCode()).Info("response_status_code")
	if err := CheckStatus(resp.StatusCode()); err != nil {
		logger.WithField("response", event).WithError(err).Error(fmt.Sprintf("failed to %s item to session", actionType))
		return nil, err
	}
	return event, nil
}

func ReplaceItem(ctx context.Context, input *model.ReplaceItemRequest) (*model.EventResponse, error) {
	logger := log.ForRequest(ctx)
	site, err := GetSiteByID(ctx, input.SiteID, input.RetailerID)
	if err != nil {
		logger.WithError(err).Error("failed to get site by id")
		return nil, err
	}
	productKeyList := make([]model.ProductKeyList, 2)
	quantityOne, quantityMinusOne := 1, -1
	productKeyList[0] = model.ProductKeyList{
		ProductKey: input.FromItem.ProductKey,
		Labelled:   input.FromItem.Labelled,
		Quantity:   &quantityMinusOne,
	}
	productKeyList[1] = model.ProductKeyList{
		ProductKey: input.ToItem.ProductKey,
		Labelled:   input.ToItem.Labelled,
		Quantity:   &quantityOne,
	}

	body := &CreateEventBody{
		ID:        input.ID,
		SessionID: input.SessionID,
		EventType: model.EventTypeShopping,
		Flagged:   input.Flagged,
		Message: model.EventMessage{
			EventSubType:   model.ReplaceActionType.ToEventSubType(),
			ProductKeyList: productKeyList,
			WebhookNotify:  input.WebhookNotify,
		},
	}
	request := buildShoppingEventRequest(ctx, site, body)

	var event *model.EventResponse
	resp, err := HttpClient().Post(ctx, request.path, request.headers, ioutil.NopCloser(bytes.NewBuffer(request.body)), &event)
	if err != nil {
		logger.WithError(err).Error(
			fmt.Sprintf("failed to post %s item event request", strings.ToLower(model.ReplaceActionType.String())))
		return nil, err
	}
	logger.WithField("resp_status_code", resp.StatusCode()).Info("response_status_code")
	if err := CheckStatus(resp.StatusCode()); err != nil {
		logger.WithField("response", event).WithError(err).Error(
			fmt.Sprintf("failed to %s item to session", strings.ToLower(model.ReplaceActionType.String())))
		return nil, err
	}
	return event, nil
}

func buildShoppingEventRequest(ctx context.Context, site *model.Site, createEventBody *CreateEventBody) ShoppingEventRequest {
	logger := log.ForRequest(ctx)
	request := ShoppingEventRequest{}

	if site.IntegrationType == model.IntegrationTypePos {
		// a tricky part, since the request body differs from the unified session service
		createEventBody.Message.Items = createEventBody.Message.ProductKeyList
		body, _ := json.Marshal(createEventBody)
		request.body = body
		// path
		request.path = Uri(RetailServicePrefix, "v2", "events")
		// headers
		headers := GenerateHeaders()
		headers["Authorization"] = middleware.GetAuthorizationHeader(ctx)
		request.headers = headers
	} else {
		body, _ := json.Marshal(createEventBody)
		request.body = body
		// path
		request.path = Uri(UnifiedSessionServicePrefix, "v1", fmt.Sprintf("sessions/%s/items", createEventBody.SessionID))
		// headers
		headers := GenerateSignatureHeaders(site.RetailerID, site.ID)
		signature, err := auth.Sign(http.MethodPost, request.path, site.RetailerID, site.ID, string(request.body), *site.SecretKey)
		if err != nil {
			logger.WithError(err).Error("failed to generate header's signature")
		}
		headers["signature"] = fmt.Sprintf("algorithm=%s, signature=%s", "HMAC-SHA256", signature)
		request.headers = headers
	}

	return request
}

func GetEventsBySessionID(ctx context.Context, sessionID, siteID, retailerID string,
	eventType *model.EventType, eventSubTypes []model.EventSubType) ([]*model.EventResponse, error) {
	// get events from the retail service first, if failed, get from the session service
	logger := log.ForRequest(ctx).WithFields(log.LogFields{
		"session_id":  sessionID,
		"site_id":     siteID,
		"retailer_id": retailerID,
	})

	queryParameter := ""
	if eventType != nil {
		queryParameter = fmt.Sprintf("?event_type=%s", *eventType)
	}
	if len(eventSubTypes) > 0 {
		for _, eventSubType := range eventSubTypes {
			if queryParameter == "" {
				queryParameter = "?"
			} else {
				queryParameter += "&"
			}
			queryParameter += fmt.Sprintf("event_sub_type=%s", eventSubType)
		}
	}

	var events []*model.EventResponse
	resp, err := HttpClient().Get(ctx,
		Uri(RetailServicePrefix, "v2", fmt.Sprintf("sessions/%s/events%s", sessionID, queryParameter)),
		GenerateHeaders(), &events)
	if err != nil {
		logger.WithError(err).
			Error("failed to get the events with the given session id from retail service")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{"response": events, "status_code": resp.StatusCode()}).WithError(err).
			Error("failed to get the events from retail service")
		return nil, err
	}

	if len(events) > 0 {
		return events, nil
	}

	// normally, if there is no timeout and network issue, err should be nil and the status should be 200
	// look up from the session service one more time
	// when the retail integration service is ready, it just needs to query from the session service
	var eventsResponse struct {
		Events []*model.EventResponse `json:"events"`
		Error  *string                `json:"error,omitempty"`
		Code   *int                   `json:"code,omitempty"`
	}
	resp, err = HttpClient().Get(ctx,
		Uri(UnifiedSessionServicePrefix, "v1",
			fmt.Sprintf("retailers/%s/sites/%s/sessions/%s/events%s",
				retailerID, siteID, sessionID, queryParameter)),
		GenerateHeaders(), &eventsResponse)
	if err != nil {
		logger.WithError(err).
			Error("failed to get the events with the given session id from session service")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		if resp.StatusCode() == http.StatusNotFound {
			return nil, nil
		}
		log.ForRequest(ctx).WithFields(log.LogFields{"response": eventsResponse, "status_code": resp.StatusCode()}).WithError(err).
			Error("failed to get the events from session service")
		return nil, err
	}

	return eventsResponse.Events, nil
}
