package model

import "gorm.io/gorm"

type TrackingStatus string

const (
	TrackingStatusProcessing = "PROCESSING"
	TrackingStatusDone       = "DONE"
)

type Tracking struct {
	gorm.Model
	TrackingStatus TrackingStatus
	VideoURLFile   string
	ZipURLFile     *string
}
