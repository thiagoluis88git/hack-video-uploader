package model

import "gorm.io/gorm"

type TrackingStatus string

const (
	TrackingStatusProcessing = "PROCESSING"
	TrackingStatusDone       = "DONE"
)

type Tracking struct {
	gorm.Model
	TrackingID        string
	TrackingStatus    TrackingStatus
	VideoURLFile      string
	ZipURLFile        *string
	ZipURLFilePresign *string
}
