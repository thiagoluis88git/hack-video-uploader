package entity

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
}
