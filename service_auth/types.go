package serviceauth

import "time"

type SyncAPIKeyPayload struct {
	Service   string     `json:"service" binding:"required"`
	Key       string     `json:"key" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
