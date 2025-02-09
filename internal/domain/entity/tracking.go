package entity

import "time"

type TrackingStatus string

const (
	TrackingStatusProcessing = "PROCESSING"
	TrackingStatusDone       = "DONE"
)

type Tracking struct {
	TrackingID     string         `json:"trackingId"`
	TrackingStatus TrackingStatus `json:"status"`
	VideoURLFile   string         `json:"videoUrl"`
	ZipURLFile     *string        `json:"zipUrl"`
	ErrorMessage   *string        `json:"errorMessage"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
