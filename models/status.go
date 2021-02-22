package models

import (
	"time"
)

type StatusType uint

const (
	UP StatusType = iota
	Down
)

type ErrorType uint

// EndpointStatus is a struct use in communication between modules
type EndpointStatus struct {
	EndpointID uint
	TimeStamp  time.Time
	Status     StatusType
	Error      ErrorType
	ErrorCode  int // Use net/http status code
}
