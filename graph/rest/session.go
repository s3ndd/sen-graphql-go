package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/s3ndd/sen-go/auth"
	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/graph/model"
	"github.com/s3ndd/sen-graphql-go/internal/middleware"
)

// SessionRequest is the body request structure for updating a session status
type SessionRequest struct {
	Paused    bool `json:"paused,omitempty"`
	Cancelled bool `json:"cancelled,omitempty"`
}

// SessionIntegrationRequest is the REST request info to session services
type SessionIntegrationRequest struct {
	path    string
	body    []byte
	headers map[string]string
}

func GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	var session model.Session
	resp, err := HttpClient().Get(ctx,
		Uri(RetailServicePrefix, "v2", fmt.Sprintf("sessions/%s", id)),
		GenerateHeaders(), &session)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("session_id", id).
			Error("failed to get the session with the given id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    session,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the session with the given id from retail service")
		return nil, err
	}
	return &session, nil
}

func GetSessionByIDs(ctx context.Context, sessionIDs []string) ([]*model.Session, error) {
	path := Uri(RetailServicePrefix, "v2", "sessions/bulk?include_unrecognised=true")
	var response map[string]*model.Session
	if err := batchQuery(ctx, path, sessionIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the session by ids from retail service")
		return nil, err
	}
	sessions := make([]*model.Session, len(sessionIDs))
	for i := range sessionIDs {
		if session, ok := response[sessionIDs[i]]; ok {
			sessions[i] = session
		}
	}

	return sessions, nil
}

func GetSessionsBySiteID(ctx context.Context, siteID string, status []model.SessionStatus) (*model.SessionConnection, error) {
	queryString := fmt.Sprintf("site_id=%s", siteID)
	for i := range status {
		queryString += fmt.Sprintf("&status=%s", status[i])
	}
	var session model.SessionConnection
	resp, err := HttpClient().Get(ctx,
		Uri(RetailServicePrefix, "v2", fmt.Sprintf("sessions?%s", queryString)),
		GenerateHeaders(), &session)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithFields(log.LogFields{
			"site_id": siteID,
			"status":  status,
		}).Error("failed to get the sessions with the given site id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    session,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the sessions with the give site id from retail service")
		return nil, err
	}
	return &session, nil
}

// GetSessionsByUserID returns the sessions of the five user id.
// This function is used to conduct whether the user is a repeat customer or a new customer.
func GetSessionsByUserID(ctx context.Context, userID string, pageIndex, pageSize int) (*model.SessionConnection, error) {
	var session model.SessionConnection
	resp, err := HttpClient().Get(ctx,
		Uri(RetailServicePrefix,
			"v2",
			fmt.Sprintf("sessions?user_id=%s&page_index=%d&page_size=%d", userID, pageIndex, pageSize)),
		GenerateHeaders(), &session)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithFields(log.LogFields{
			"user_id":    userID,
			"page_index": pageIndex,
			"page_size":  pageSize,
		}).Error("failed to get the sessions with the given site id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    session,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the sessions with the give site id from retail service")
		return nil, err
	}
	return &session, nil
}

func GetUnifiedSessionByIDs(ctx context.Context, sessionIDs []string) ([]*model.Session, error) {
	path := Uri(UnifiedSessionServicePrefix, "v1", "sessions/bulk?include_unrecognised=true")
	var response map[string]*model.Session
	if err := batchQuery(ctx, path, sessionIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the session by ids from unified session service")
		return nil, err
	}
	sessions := make([]*model.Session, len(sessionIDs))
	for i := range sessionIDs {
		if session, ok := response[sessionIDs[i]]; ok {
			sessions[i] = session
		}
	}

	return sessions, nil
}

// UpdateSessionById returns the session of the given session id.
// This function is used to update the status of a session and routes the request respective service
func UpdateSessionById(ctx context.Context, input *model.UpdateSessionStatusRequest) (*model.Session, error) {
	logger := log.ForRequest(ctx).WithField("input", input)
	request := updateSessionStatusHelper(ctx, input)
	logger = logger.WithField("request", request)
	if request == nil {
		err := fmt.Errorf("invalid update session status request")
		logger.WithError(err).Error(err.Error())
		return nil, err
	}
	var session *model.Session
	resp, err := HttpClient().Post(ctx, request.path, request.headers, ioutil.NopCloser(bytes.NewBuffer(request.body)), &session)
	if err != nil {
		logger.WithError(err).Error("failed to post session status update request")
		return nil, err
	}
	logger.WithField("resp_status_code", resp.StatusCode()).Info("response_status_code")
	if err := CheckStatus(resp.StatusCode()); err != nil {
		logger.WithField("response", session).WithError(err).Error("failed to update session status with the given id and status")
		return nil, err
	}
	return session, nil
}

// A helper function to construct a REST request's headers, path and body
func updateSessionStatusHelper(ctx context.Context, input *model.UpdateSessionStatusRequest) *SessionIntegrationRequest {
	logger := log.ForRequest(ctx).WithField("input", input)
	site, err := GetSiteByID(ctx, input.SiteID, input.RetailerID)
	if err != nil {
		logger.WithError(err).Error("failed to get site by id")
		return nil
	}

	useRetailService := site.IntegrationType == model.IntegrationTypePos
	request := &SessionIntegrationRequest{}
	defer func(useRetailService bool, request *SessionIntegrationRequest) {
		if useRetailService {
			headers := GenerateHeaders()
			headers["Authorization"] = middleware.GetAuthorizationHeader(ctx)
			request.headers = headers
		} else {
			headers := GenerateSignatureHeaders(input.RetailerID, input.SiteID)
			signature, err := auth.Sign(http.MethodPost, request.path, input.RetailerID, input.SiteID, string(request.body), *site.SecretKey)
			if err != nil {
				logger.WithError(err).Error("failed to generate header's signature")
			}
			headers["signature"] = fmt.Sprintf("algorithm=%s, signature=%s", "HMAC-SHA256", signature)
			request.headers = headers
		}
	}(useRetailService, request)

	switch input.Status {
	case model.SessionStatusPaused:
		request.body, _ = json.Marshal(&SessionRequest{
			Paused: true,
		})
		if useRetailService {
			request.path = Uri(RetailServicePrefix, "v2", fmt.Sprintf("sessions/%s/paused", input.SessionID))
		} else {
			request.path = Uri(UnifiedSessionServicePrefix, "v1", fmt.Sprintf("sessions/%s/pause", input.SessionID))
		}
	case model.SessionStatusShopping:
		request.body, _ = json.Marshal(&SessionRequest{
			Paused: false,
		})
		if useRetailService {
			request.path = Uri(RetailServicePrefix, "v2", fmt.Sprintf("sessions/%s/paused", input.SessionID))
		} else {
			request.path = Uri(UnifiedSessionServicePrefix, "v1", fmt.Sprintf("sessions/%s/resume", input.SessionID))
		}
	case model.SessionStatusCancelled:
		request.body, _ = json.Marshal(&SessionRequest{
			Cancelled: true,
		})
		if useRetailService {
			request.path = Uri(RetailServicePrefix, "v2", fmt.Sprintf("sessions/%s/cancelled", input.SessionID))
		} else {
			request.path = Uri(UnifiedSessionServicePrefix, "v1", fmt.Sprintf("sessions/%s/cancel", input.SessionID))
		}
	default:
		return nil
	}

	return request
}
