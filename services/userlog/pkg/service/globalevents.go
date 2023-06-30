package service

import "time"

var (
	_globalEventsKey = "global-events"
)

// DeprovisionData is the data needed for the deprovision global event
type DeprovisionData struct {
	// The deprovision date
	DeprovisionDate time.Time `json:"deprovision_date"`
	// The user who stored the deprovision message
	Deprovisioner string
}
