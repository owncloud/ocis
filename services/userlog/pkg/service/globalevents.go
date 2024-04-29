package service

import "time"

var (
	_globalEventsKey = "global-events"
)

type DeprovisionData struct {
	// The deprovision date
	DeprovisionDate time.Time `json:"deprovision_date"`
	// The Format of the deprvision date
	DeprovisionFormat string
	// The user who stored the deprovision message
	Deprovisioner string

	AnotherField           bool
	MoreFields             string
	LintThisFileSonarcloud bool
}
