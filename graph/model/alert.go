package model

import "time"

// AlertStatus is the type of alert status
type AlertStatus string

const (
	// OpenAlertStatus is the initial status that means the alert has to be acknowledged, resolved or self solved.
	OpenAlertStatus AlertStatus = "OPEN"
	// AcknowledgedAlertStatus is the status that means the alert has been claimed by staff.
	AcknowledgedAlertStatus AlertStatus = "ACKNOWLEDGED"
	// ResolvedAlertStatus is the final status that means the alert has been resolved.
	ResolvedAlertStatus AlertStatus = "RESOLVED"
	// SelfSolvedAlertStatus is the final status that means the alert has been solved by the user.
	SelfSolvedAlertStatus AlertStatus = "SELF_SOLVED"
	// NaAlertStatus is not an actual status in the database, which is only used in the code.
	NaAlertStatus AlertStatus = "NA"
)

// AlertType is the type of alert
type AlertType string

const (
	HelpAlertType          AlertType = "HELP"
	LatchAlertType         AlertType = "LATCH"
	POSErrorAlertType      AlertType = "POS_ERROR"
	MultipleItemsAlertType AlertType = "MULTIPLE_ITEMS"
	MissedLabelAlertType   AlertType = "MISSED_LABEL"
)

// Alert defines the alert body.
type Alert struct {
	ID             string      `json:"id"`
	SiteID         string      `json:"site_id"`
	CartID         string      `json:"cart_id"`
	CartQRCode     string      `json:"cart_qr_code"`
	SessionID      *string     `json:"session_id"`
	Status         AlertStatus `json:"status"`
	Type           AlertType   `json:"type"`
	Responder      *string     `json:"responder"`
	Message        *string     `json:"message"`
	TriggeredAt    time.Time   `json:"triggered_at"`
	AcknowledgedAt *time.Time  `json:"acknowledged_at"`
	ResolvedAt     *time.Time  `json:"resolved_at"`
	Created        time.Time   `json:"created"`
	Updated        time.Time   `json:"updated"`
	Deleted        bool        `json:"deleted"`
}

type AlertConnection struct {
	Alerts []*Alert `json:"alerts"`
}
