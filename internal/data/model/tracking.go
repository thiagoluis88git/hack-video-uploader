package model

import "gorm.io/gorm"

type TrackingStatus string

const (
	TrackingStatusProcessing = "PROCESSING"
	TrackingStatusDone       = "DONE"
	TrackingStatusError      = "ERROR"
)

type Tracking struct {
	gorm.Model
	TrackingID        string
	CPF               string
	TrackingStatus    TrackingStatus
	VideoURLFile      string
	ZipURLFile        *string
	ZipURLFilePresign *string
	ErrorMessage      *string
}
