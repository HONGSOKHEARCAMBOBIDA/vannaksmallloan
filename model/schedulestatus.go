package model

type ScheduelStatus int

const (
	UnknownScheduleStatus ScheduelStatus = iota
	PENDING
	PAID
	PARTIAL
	OVERDUE
	CANCELLED
)
