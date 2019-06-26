package models

import (
	"time"
)

const (
	RunStatusReadyForCopyTestpack   = 0
	RunStatusCopyTestpackInProgress = 1
	RunStatusCopyTestpackFailed     = 2
	RunStatusReadyForNPMInstall     = 3
	RunStatusNPMInstallProgress     = 4
	RunStatusNPMInstallFailed       = 5
	RunStatusReadyForCafeThread     = 6
	RunStatusCafeThreadFailed       = 7
	RunStatusCafeThreadProgress     = 8
	RunStatusCafeThreadFinTimeout   = 9
	RunStatusCafeThreadFin          = 10 // testcafe process finished with some exit code (zero or not).
)

type Run struct {
	ID              int `storm:"increment"`
	SessionID       int
	Status          uint8
	StartTime       time.Time
	Port            int
	DeviceOwnerName string
	ExitCode        string
	StdOut          []byte
}
