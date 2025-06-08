package model

import "time"

// SiteStatus indicates the status of a site
type SiteStatus string

const (
	// ActiveSiteStatus indicates the site is active
	ActiveSiteStatus SiteStatus = "ACTIVE"
)

// SiteRegion indicates the region a site is located in
type SiteRegion string

const (
	// AsiaSiteRegion indicates the site is located in Asia
	AsiaSiteRegion SiteRegion = "ASIA"
	// AustraliaSiteRegion indicates the site is located in Australia
	AustraliaSiteRegion SiteRegion = "AUSTRALIA"
	// EuropeSiteRegion indicates the site is located in Europe
	EuropeSiteRegion SiteRegion = "EUROPE"
	// USSiteRegion indicates the site is located in US
	USSiteRegion SiteRegion = "US"
)

// Site contains information of a site
type Site struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Status                SiteStatus      `json:"status"`
	Region                SiteRegion      `json:"region"`
	Currency              string          `json:"currency"`
	WorkflowType          string          `json:"workflow_type"`
	IntegrationType       IntegrationType `json:"integration_type"`
	SecretKey             *string         `json:"secret_key,omitempty"`
	AlertNotificationType string          `json:"alert_notification_type"`
	AlertNotificationURL  *string         `json:"alert_notification_url,omitempty"`
	RetailerID            string          `json:"retailer_id"`
	Retailer              *Retailer       `json:"retailer"`
	LogoURL               string          `json:"logo_url,omitempty"`
	Created               time.Time       `json:"created"`
	Updated               time.Time       `json:"updated"`
	Deleted               bool            `json:"deleted"`
}
