package model

import "time"

type DataLogger struct {
	StartTime    time.Time
	EndTime      time.Time
	ElapsedNano  int64
	Error        error
	Body         string
	Where        string
	Operation    string
	BodyResponse string
	IsError      bool
}
