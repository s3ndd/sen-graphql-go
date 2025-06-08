package model

import "time"

// RetailerStatus indicates the status of a retailer
type RetailerStatus string

const (
	// ActiveRetailerStatus indicates retailer is an active customer
	ActiveRetailerStatus RetailerStatus = "ACTIVE"
	// PrePilotRetailerStatus indicates retailer is a prepilot customer
	PrePilotRetailerStatus RetailerStatus = "PREPILOT"
	// PilotRetailerStatus indicates retailer is a pilot customer
	PilotRetailerStatus RetailerStatus = "PILOT"
	// DemoRetailerStatus indicates retailer is a demo customer
	DemoRetailerStatus RetailerStatus = "DEMO"
	// TrialRetailerStatus indicates retailer is a trial customer
	TrialRetailerStatus RetailerStatus = "TRIAL"
)

// Retailer contains information of a retailer
type Retailer struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Status  RetailerStatus `json:"status"`
	PdfURL  string         `json:"pdf_url,omitempty"`
	Created time.Time      `json:"created"`
	Updated time.Time      `json:"updated"`
	Deleted bool           `json:"deleted"`
}
