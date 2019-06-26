package models

import (
	"time"
)

// Testpack model
type Session struct {
	ID          int `storm:"increment"`
	TestpackID  int
	CreatedTime time.Time
}
