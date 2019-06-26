package models

import (
	"time"
)

const (
	TPStatusReadyForUnzip   = 0
	TPStatusUnzipInProgress = 1
	TPStatusUnzipFailed     = 2
	// TPStatusReadyForInit    = 3
	// TPStatusInitProgress    = 4
	// TPStatusInitFailed      = 5
	TPStatusReadyForRunning = 3
)

type Testpack struct {
	ID          int `storm:"increment"`
	Status      uint8
	UploadTime  time.Time
	Zip         []byte
	ZipHash     []byte
	InitFailOut []byte
}
