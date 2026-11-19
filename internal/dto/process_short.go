package dto

import "time"

type ProcessShort struct {
	PID         int64
	ProcessName string
	Path        string
	CreatedAt   time.Time
}
