package rest

import (
	"context"
	"fmt"

	"github.com/s3ndd/sen-go/log"

	"github.com/s3ndd/sen-graphql-go/graph/model"
)

func GetAlertByID(ctx context.Context, id string) (*model.Alert, error) {
	var alert model.Alert
	resp, err := HttpClient().Get(ctx,
		Uri(AlertNotificationServicePrefix, "v1", fmt.Sprintf("alerts/%s", id)),
		GenerateHeaders(), &alert)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("alert_id", id).
			Error("failed to get the alert with the given id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    alert,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the alert from alert notification service")
		return nil, err
	}
	return &alert, nil
}

func GetAlertByIDs(ctx context.Context, alertIDs []string) ([]*model.Alert, error) {
	path := Uri(AlertNotificationServicePrefix, "v1", "alerts/bulk")
	var response map[string]*model.Alert
	if err := batchQuery(ctx, path, alertIDs, &response); err != nil {
		log.ForRequest(ctx).WithError(err).
			Error("failed to get the alerts by ids from cart service")
		return nil, err
	}
	alerts := make([]*model.Alert, len(response))
	for i := range alertIDs {
		if alert, ok := response[alertIDs[i]]; ok {
			alerts[i] = alert
		}
	}

	return alerts, nil
}

func GetAlerts(ctx context.Context, siteID, sessionID, cartID *string,
	status []model.AlertStatus, types []model.AlertType) (*model.AlertConnection, error) {
	var alerts model.AlertConnection
	queryString := ""
	if siteID != nil {
		queryString = fmt.Sprintf("?site_id=%s", *siteID)
	}
	if sessionID != nil {
		symbol := "?"
		if queryString != "" {
			symbol = "&"
		}
		queryString += fmt.Sprintf("%ssession_id=%s", symbol, *sessionID)
	}
	if cartID != nil {
		symbol := "?"
		if queryString != "" {
			symbol = "&"
		}
		queryString += fmt.Sprintf("%scart_id=%s", symbol, *cartID)
	}

	statusString, typesString := "", ""
	for i := range status {
		statusString += string(status[i])
		if i < len(status)-1 {
			statusString += ","
		}
	}
	if statusString != "" {
		symbol := "?"
		if queryString != "" {
			symbol = "&"
		}
		queryString += fmt.Sprintf("%sstatus=%s", symbol, statusString)
	}
	for i := range types {
		typesString += string(types[i])
		if i < len(types)-1 {
			typesString += ","
		}
	}
	if typesString != "" {
		symbol := "?"
		if queryString != "" {
			symbol = "&"
		}

		queryString += fmt.Sprintf("%stypes=%s", symbol, typesString)
	}

	resp, err := HttpClient().Get(ctx,
		Uri(AlertNotificationServicePrefix, "v1",
			fmt.Sprintf("alerts%s", queryString)),
		GenerateHeaders(), &alerts)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("query_string", queryString).
			Error("failed to get the alerts with the given parameters")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    alerts,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the alerts from alert notification service")
		return nil, err
	}
	return &alerts, nil
}

func GetAlertsBySessionID(ctx context.Context, sessionID string, status []model.AlertStatus, types []model.AlertType) ([]*model.Alert, error) {
	var alerts struct {
		Alerts []*model.Alert `json:"alerts"`
	}
	statusString, typesString := "", ""
	for i := range status {
		statusString += string(status[i])
		if i < len(status)-1 {
			statusString += ","
		}
	}
	if statusString != "" {
		statusString = fmt.Sprintf("&status=%s", statusString)
	}
	for i := range types {
		typesString += string(types[i])
		if i < len(types)-1 {
			typesString += ","
		}
	}
	if typesString != "" {
		typesString = fmt.Sprintf("&types=%s", typesString)
	}

	resp, err := HttpClient().Get(ctx,
		Uri(AlertNotificationServicePrefix, "v1",
			fmt.Sprintf("alerts?session_id=%s%s%s", sessionID, statusString, typesString)),
		GenerateHeaders(), &alerts)
	if err != nil {
		log.ForRequest(ctx).WithError(err).WithField("session_id", sessionID).
			Error("failed to get the alerts with the given session id")
		return nil, err
	}

	if err := CheckStatus(resp.StatusCode()); err != nil {
		log.ForRequest(ctx).WithFields(log.LogFields{
			"response":    alerts,
			"status_code": resp.StatusCode(),
		}).WithError(err).Error("failed to get the alerts from alert notification service")
		return nil, err
	}
	return alerts.Alerts, nil
}

func GetAlertsBySiteID(ctx context.Context, siteID string) ([]*model.Alert, error) {
	return nil, nil
}
