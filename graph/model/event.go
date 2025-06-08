package model

import "time"

type EventResponse struct {
	ID        string       `json:"id"`
	SessionID string       `json:"session_id"`
	EventType EventType    `json:"event_type"`
	Flagged   bool         `json:"flagged"`
	Skipped   bool         `json:"skipped"`
	Message   EventMessage `json:"message"`
	Timestamp *time.Time   `json:"timestamp,omitempty"`
	Created   time.Time    `json:"created"`
	Updated   time.Time    `json:"updated"`
}

type Event struct {
	ID             string           `json:"id"`
	SessionID      string           `json:"sessionID"`
	EventType      EventType        `json:"eventType"`
	EventSubType   EventSubType     `json:"eventSubType"`
	ProductKeyList []ProductKeyList `json:"productKeyList"`
	Flagged        bool             `json:"flagged"`
	Skipped        bool             `json:"skipped"`
	Created        time.Time        `json:"created"`
	Updated        time.Time        `json:"updated"`
}

type EventMessage struct {
	EventSubType   EventSubType     `json:"event_sub_type"`
	WebhookNotify  bool             `json:"webhook_notify"`
	ProductKeyList []ProductKeyList `json:"product_key_list,omitempty"`
	// Items filed is used to send request to the retail service
	Items []ProductKeyList `json:"items,omitempty"`
}

type ActionType string

const (
	AddActionType     ActionType = "ADD"
	RemoveActionType  ActionType = "REMOVE"
	ReplaceActionType ActionType = "REPLACE"
)

func (a ActionType) String() string {
	return string(a)
}

func (a ActionType) ToEventSubType() EventSubType {
	return EventSubType(a)
}

type ProductKeyList struct {
	ProductKey       string   `json:"product_key" binding:"required"`
	Labelled         bool     `json:"labelled"`
	CandidateClasses []string `json:"candidate_classes,omitempty"`
	BarcodeReader    *string  `json:"barcode_reader_product_key,omitempty"`
	Discount         *string  `json:"discount,omitempty"`
	Quantity         *int     `json:"quantity"`
}

type ItemsRequest struct {
	ID            string           `json:"id"`
	SessionID     string           `json:"sessionID"`
	SiteID        string           `json:"siteID"`
	RetailerID    string           `json:"retailerID"`
	Flagged       bool             `json:"flagged"`
	Skipped       bool             `json:"skipped"`
	Items         []ProductKeyList `json:"item"`
	WebhookNotify bool             `json:"webhook_notify"`
}

type ReplaceItemRequest struct {
	ID            string         `json:"id"`
	SessionID     string         `json:"sessionID"`
	SiteID        string         `json:"siteID"`
	RetailerID    string         `json:"retailerID"`
	Flagged       bool           `json:"flagged"`
	Skipped       bool           `json:"skipped"`
	FromItem      ProductKeyList `json:"fromItem"`
	ToItem        ProductKeyList `json:"toItem"`
	WebhookNotify bool           `json:"webhook_notify"`
}
