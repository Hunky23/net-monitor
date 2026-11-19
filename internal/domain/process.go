package domain

import "time"

const (
	ProcessStatusActive   = "active"
	ProcessStatusFinished = "finished"
)

// Типы процессов.
const (
	ProcessTypeCore    = "core"
	ProcessTypeCommon  = "common"
	ProcessTypeUWP     = "uwp"
	ProcessTypeService = "service"
)

type Process struct {
	ID          int64
	PID         int64
	ProcessName string
	DisplayName string
	Path        string
	Status      string
	Type        string
	Icon        string
	CreatedAt   time.Time
}
